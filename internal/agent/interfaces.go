package agent

// Planner analyzes the repo context and creates a test plan.
type Planner interface {
	Plan(repoContext string) (string, error)
}

// Coder generates code based on the plan and target file.
type Coder interface {
	GenerateCode(plan string, fileContent string) (string, error)
}

// Executor runs the generated code and returns logs.
type Executor interface {
	Execute(code string) (string, error)
}

// Agent orchestrates the components.
type Agent struct {
	Planner  Planner
	Coder    Coder
	Executor Executor
}

func NewAgent(p Planner, c Coder, e Executor) *Agent {
	return &Agent{
		Planner:  p,
		Coder:    c,
		Executor: e,
	}
}

func (a *Agent) Run(repoContext, fileContent string) error {
	// 1. Plan
	plan, err := a.Planner.Plan(repoContext)
	if err != nil {
		return err
	}

	// 2. Code
	code, err := a.Coder.GenerateCode(plan, fileContent)
	if err != nil {
		return err
	}

	// 3. Execute
	_, err = a.Executor.Execute(code)
	if err != nil {
		return err
	}
	
	return nil
}
