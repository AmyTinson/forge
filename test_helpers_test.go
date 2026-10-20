package main

import (
	"context"
	"os"
	"testing"
)

type registryTestWorkflow struct {
	ctx  context.Context
	args []string
}

type fakeLLMProvider struct {
	prompts             []string
	outputs             []string
	failAt              int
	failWith            error
	waitForCancellation bool
	started             chan struct{}
}

func (p *fakeLLMProvider) Generate(ctx context.Context, prompt string) (string, error) {
	p.prompts = append(p.prompts, prompt)
	call := len(p.prompts) - 1
	if p.failWith != nil && call == p.failAt {
		return "", p.failWith
	}
	if p.waitForCancellation {
		close(p.started)
		<-ctx.Done()
		return "", ctx.Err()
	}
	return p.outputs[call], nil
}

func (w *registryTestWorkflow) Name() string {
	return "test"
}

func (w *registryTestWorkflow) Run(ctx context.Context) error {
	w.ctx = ctx
	return nil
}

func setWorkingDirectory(t *testing.T, dir string) {
	t.Helper()
	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(original); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})
}
