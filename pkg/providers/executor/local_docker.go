package executor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type LocalDockerExecutor struct {
	Config ExecutorConfig
}

func NewLocalDockerExecutor(cfg ExecutorConfig) *LocalDockerExecutor {
	// Apply defaults if not set
	if cfg.WorkDir == "" {
		cfg.WorkDir = "/app"
	}
	if cfg.TestFilePattern == "" {
		cfg.TestFilePattern = "generated_test.go"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 300
	}
	if len(cfg.Command) == 0 {
		cfg.Command = []string{"go", "test", "-v", "./..."}
	}

	return &LocalDockerExecutor{Config: cfg}
}

func (l *LocalDockerExecutor) Execute(code string) (string, error) {
	timeout := time.Duration(l.Config.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create Docker client using default socket
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create docker client: %w", err)
	}
	defer cli.Close()

	fmt.Printf("[Executor] Connected to local Docker\n")

	// Write code to a temporary file that will be mounted
	tempDir, err := os.MkdirTemp("", "localsprite-test-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, l.Config.TestFilePattern)
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		return "", fmt.Errorf("failed to write test file: %w", err)
	}

	// Ensure the image is available
	fmt.Printf("[Executor] Pulling image %s (if needed)...\n", l.Config.Image)
	pullOut, err := cli.ImagePull(ctx, l.Config.Image, image.PullOptions{})
	if err != nil {
		// Image might already exist locally, continue anyway
		fmt.Printf("[Executor] Image pull skipped: %v\n", err)
	} else {
		io.Copy(io.Discard, pullOut)
		pullOut.Close()
	}

	// Create container configuration
	containerConfig := &container.Config{
		Image:      l.Config.Image,
		Cmd:        l.Config.Command,
		WorkingDir: l.Config.WorkDir,
		Tty:        false,
	}

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: tempDir,
				Target: l.Config.WorkDir,
			},
		},
		AutoRemove: false,
	}

	// Create the container
	fmt.Printf("[Executor] Creating container with command: %v\n", l.Config.Command)
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}
	containerID := resp.ID

	// Ensure cleanup
	defer func() {
		removeCtx, removeCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer removeCancel()
		cli.ContainerRemove(removeCtx, containerID, container.RemoveOptions{Force: true})
	}()

	// Start the container
	fmt.Printf("[Executor] Starting container %s...\n", containerID[:12])
	if err := cli.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	// Wait for container to finish
	statusCh, errCh := cli.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", fmt.Errorf("error waiting for container: %w", err)
		}
	case status := <-statusCh:
		fmt.Printf("[Executor] Container exited with code %d\n", status.StatusCode)
	}

	// Get container logs
	logOptions := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	}
	logs, err := cli.ContainerLogs(ctx, containerID, logOptions)
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}
	defer logs.Close()

	// Demultiplex stdout/stderr
	var stdout, stderr bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdout, &stderr, logs); err != nil {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}

	output := stdout.String()
	if stderr.Len() > 0 {
		output += "\n--- STDERR ---\n" + stderr.String()
	}

	return output, nil
}
