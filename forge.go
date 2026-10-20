package main

import (
	"context"
	"fmt"
	"os"
)

func newDefaultWorkflowRegistry() (*WorkflowRegistry, error) {
	registry := NewWorkflowRegistry()
	if err := registry.Register(bugTriageWorkflowName, newBugTriageWorkflow); err != nil {
		return nil, err
	}
	if err := registry.Register(
		workflowDefinitionCommandName,
		NewWorkflowDefinitionFactory(func() (LLMProvider, error) {
			return HTTPGenerationProviderFromEnv()
		}),
	); err != nil {
		return nil, err
	}
	return registry, nil
}

func runRead(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	if _, err := os.Stdout.Write(file); err != nil {
		return fmt.Errorf("write to stdout: %w", err)
	}
	return nil
}

func forge(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: forge <command> [arguments]")
	}

	switch args[0] {
	case "read":
		if len(args) != 2 || args[1] == "" {
			return fmt.Errorf("usage: forge read <file>")
		}
		return runRead(args[1])
	case "run":
		registry, err := newDefaultWorkflowRegistry()
		if err != nil {
			return fmt.Errorf("register workflows: %w", err)
		}
		return runCommand(registry, args[1:])
	case "note":
		return noteCommand(args[1:])
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func runCommand(registry *WorkflowRegistry, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: forge run <workflow> [options]")
	}

	return registry.Run(context.Background(), args[0], args[1:])
}

func noteCommand(args []string) error {
	if len(args) == 0 || args[0] != "add" {
		if len(args) == 0 {
			return fmt.Errorf("usage: forge note add \"your note\"")
		}
		return fmt.Errorf("unknown note command: %s", args[0])
	}

	if len(args) != 2 || args[1] == "" {
		return fmt.Errorf("usage: forge note add \"your note\"")
	}
	if err := addNote(args[1]); err != nil {
		return fmt.Errorf("add note: %w", err)
	}
	fmt.Println("note added:", args[1])
	return nil
}
