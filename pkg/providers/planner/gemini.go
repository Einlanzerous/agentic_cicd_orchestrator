package planner

import (
	"fmt"
)

type GeminiPlanner struct {
	Model string
}

func NewGeminiPlanner(model string) *GeminiPlanner {
	return &GeminiPlanner{Model: model}
}

func (g *GeminiPlanner) Plan(repoContext string) (string, error) {
	fmt.Printf("[Planner] Gemini (%s) is analyzing repository context...\n", g.Model)
	// Integration with Google GenAI SDK would go here.
	return "Test Plan: Cover edge cases in auth module", nil
}