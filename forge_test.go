package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type registryTestWorkflow struct {
	ctx  context.Context
	args []string
}

func (w *registryTestWorkflow) Name() string {
	return "test"
}

func (w *registryTestWorkflow) Run(ctx context.Context) error {
	w.ctx = ctx
	return nil
}

func TestForgeRejectsInvalidArguments(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "missing command"},
		{name: "missing read path", args: []string{"read"}},
		{name: "extra read argument", args: []string{"read", "input.md", "extra"}},
		{name: "unknown command", args: []string{"unknown"}},
		{name: "missing workflow", args: []string{"run"}},
		{name: "unknown workflow", args: []string{"run", "unknown"}},
		{name: "missing workflow input", args: []string{"run", "bug-triage"}},
		{name: "unknown workflow option", args: []string{"run", "bug-triage", "--unknown"}},
		{name: "unexpected workflow argument", args: []string{"run", "bug-triage", "--input", "report.md", "extra"}},
		{name: "missing note operation", args: []string{"note"}},
		{name: "missing note value", args: []string{"note", "add"}},
		{name: "extra note argument", args: []string{"note", "add", "one", "two"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := forge(test.args); err == nil {
				t.Fatalf("forge(%v) succeeded, want an error", test.args)
			}
		})
	}
}

func TestWorkflowRegistryRunsRegisteredWorkflow(t *testing.T) {
	workflow := &registryTestWorkflow{}
	registry := NewWorkflowRegistry()
	err := registry.Register("test", func(args []string) (Workflow, error) {
		workflow.args = args
		return workflow, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	if err := registry.Run(ctx, "test", []string{"--value", "expected"}); err != nil {
		t.Fatalf("registry.Run: %v", err)
	}
	if workflow.ctx != ctx {
		t.Fatal("workflow did not receive the run context")
	}
	if strings.Join(workflow.args, " ") != "--value expected" {
		t.Fatalf("workflow arguments = %v, want [--value expected]", workflow.args)
	}
}

func TestBugTriageCommandPersistsDistinctRunOutputs(t *testing.T) {
	dir := t.TempDir()
	setWorkingDirectory(t, dir)

	inputPath := filepath.Join(dir, "report.md")
	if err := os.WriteFile(inputPath, []byte("A reproducible bug report."), 0644); err != nil {
		t.Fatal(err)
	}

	args := []string{"run", "bug-triage", "--input", inputPath}
	for range 2 {
		if err := forge(args); err != nil {
			t.Fatalf("forge(%v): %v", args, err)
		}
	}

	runPaths, err := filepath.Glob(filepath.Join("data", "runs", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(runPaths) != 2 {
		t.Fatalf("got %d run records, want 2", len(runPaths))
	}

	artifactNames := make(map[string]bool)
	for _, runPath := range runPaths {
		runData, err := os.ReadFile(runPath)
		if err != nil {
			t.Fatal(err)
		}
		var run Run
		if err := json.Unmarshal(runData, &run); err != nil {
			t.Fatalf("decode run record %q: %v", runPath, err)
		}
		if artifactNames[run.ArtifactFileName] {
			t.Fatalf("multiple runs refer to artifact %q", run.ArtifactFileName)
		}
		artifactNames[run.ArtifactFileName] = true

		artifactPath := filepath.Join("data", "artifacts", run.ArtifactFileName)
		artifact, err := os.ReadFile(artifactPath)
		if err != nil {
			t.Fatalf("read artifact %q: %v", artifactPath, err)
		}
		if !strings.Contains(string(artifact), "A reproducible bug report.") {
			t.Errorf("artifact %q does not contain the input report", artifactPath)
		}
	}
}

func TestWriteNewFileDoesNotOverwrite(t *testing.T) {
	path := filepath.Join(t.TempDir(), "existing.txt")
	if err := os.WriteFile(path, []byte("original"), 0644); err != nil {
		t.Fatal(err)
	}

	err := writeNewFile(path, []byte("replacement"))
	if !errors.Is(err, os.ErrExist) {
		t.Fatalf("writeNewFile error = %v, want os.ErrExist", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "original" {
		t.Fatalf("file content = %q, want original content", data)
	}
}

func TestNoteErrorsReachCommandBoundary(t *testing.T) {
	dir := t.TempDir()
	setWorkingDirectory(t, dir)
	if err := os.WriteFile("notes", []byte("not a directory"), 0644); err != nil {
		t.Fatal(err)
	}

	err := forge([]string{"note", "add", "test note"})
	if err == nil || !strings.Contains(err.Error(), "add note") {
		t.Fatalf("forge error = %v, want an add-note error", err)
	}
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
