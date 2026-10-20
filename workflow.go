package main

import (
	"context"
	"fmt"
)

type Workflow interface {
	Name() string
	Run(ctx context.Context) error
}

type WorkflowFactory func(args []string) (Workflow, error)

type WorkflowRegistry struct {
	factories map[string]WorkflowFactory
}

func NewWorkflowRegistry() *WorkflowRegistry {
	return &WorkflowRegistry{factories: make(map[string]WorkflowFactory)}
}

func (r *WorkflowRegistry) Register(name string, factory WorkflowFactory) error {
	if name == "" {
		return fmt.Errorf("workflow name cannot be empty")
	}
	if factory == nil {
		return fmt.Errorf("workflow factory for %q cannot be nil", name)
	}
	if _, exists := r.factories[name]; exists {
		return fmt.Errorf("workflow %q is already registered", name)
	}

	r.factories[name] = factory
	return nil
}

func (r *WorkflowRegistry) Run(ctx context.Context, name string, args []string) error {
	factory, exists := r.factories[name]
	if !exists {
		return fmt.Errorf("unknown workflow: %s", name)
	}

	workflow, err := factory(args)
	if err != nil {
		return err
	}
	if workflow == nil {
		return fmt.Errorf("workflow factory %q returned nil", name)
	}
	if workflow.Name() != name {
		return fmt.Errorf("workflow factory %q returned workflow %q", name, workflow.Name())
	}
	return workflow.Run(ctx)
}
