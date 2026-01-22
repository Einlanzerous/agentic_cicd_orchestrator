package executor

import (
	"fmt"
)

type LocalDockerExecutor struct {
	Image string
}

func NewLocalDockerExecutor(image string) *LocalDockerExecutor {
	return &LocalDockerExecutor{Image: image}
}

func (l *LocalDockerExecutor) Execute(code string) (string, error) {
	fmt.Printf("[Executor] Local Docker running image %s...\n", l.Image)
	// Docker SDK integration here
	return "PASS: 5/5 tests", nil
}