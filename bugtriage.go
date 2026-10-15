package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AmyTinson/forge/internal/triage"

	"github.com/google/uuid"
)

const bugTriageWorkflowName = "bug-triage"

type Run struct {
	ID               string    `json:"id"`
	Workflow         string    `json:"workflow"`
	InputFileName    string    `json:"inputFileName"`
	ArtifactFileName string    `json:"artifactFileName"`
	CreatedAt        time.Time `json:"createdAt"`
	Error            string    `json:"error,omitempty"`
}

type BugTriageWorkflow struct {
	InputFile string
}

var _ Workflow = BugTriageWorkflow{}

func (w BugTriageWorkflow) Name() string {
	return bugTriageWorkflowName
}

func init() {
	if err := workflowRegistry.Register(bugTriageWorkflowName, newBugTriageWorkflow); err != nil {
		panic(err)
	}
}

func newBugTriageWorkflow(args []string) (Workflow, error) {
	flags := flag.NewFlagSet(bugTriageWorkflowName, flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	inputFile := flags.String("input", "", "input bug report file")
	if err := flags.Parse(args); err != nil {
		return nil, fmt.Errorf("parse bug-triage options: %w", err)
	}
	if flags.NArg() != 0 || *inputFile == "" {
		return nil, fmt.Errorf("usage: forge run bug-triage --input <file>")
	}
	return BugTriageWorkflow{InputFile: *inputFile}, nil
}

func (w BugTriageWorkflow) Run(ctx context.Context) error {
	inputFile := w.InputFile
	bugReportFile, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}

	artifactDir := filepath.Join("data", "artifacts")
	runDir := filepath.Join("data", "runs")
	if err := os.MkdirAll(artifactDir, 0755); err != nil {
		return fmt.Errorf("create artifact directory: %w", err)
	}
	if err := os.MkdirAll(runDir, 0755); err != nil {
		return fmt.Errorf("create run directory: %w", err)
	}

	triageData := triage.TriageData{
		FileName:  filepath.Base(inputFile),
		RawReport: string(bugReportFile),
	}
	renderedTriage, err := triage.RenderTriage(triageData)
	if err != nil {
		return fmt.Errorf("render triage: %w", err)
	}

	inputFileName := filepath.Base(inputFile)
	inputFileExt := filepath.Ext(inputFileName)
	runID := uuid.NewString()
	run := Run{
		ID:               runID,
		Workflow:         w.Name(),
		InputFileName:    inputFileName,
		ArtifactFileName: fmt.Sprintf("%s-triage-%s.md", strings.TrimSuffix(inputFileName, inputFileExt), runID),
		CreatedAt:        time.Now(),
	}
	runData, err := json.Marshal(run)
	if err != nil {
		return fmt.Errorf("marshal run record: %w", err)
	}

	artifactPath := filepath.Join(artifactDir, run.ArtifactFileName)
	if err := writeNewFile(artifactPath, []byte(renderedTriage+"\n")); err != nil {
		return fmt.Errorf("write triage artifact: %w", err)
	}

	runPath := filepath.Join(runDir, run.ID+".json")
	if err := writeNewFile(runPath, runData); err != nil {
		if cleanupErr := os.Remove(artifactPath); cleanupErr != nil {
			return errors.Join(
				fmt.Errorf("write run record: %w", err),
				fmt.Errorf("remove incomplete artifact %q: %w", artifactPath, cleanupErr),
			)
		}
		return fmt.Errorf("write run record: %w", err)
	}

	fmt.Println("triage written:", artifactPath)
	return nil
}

func writeNewFile(path string, data []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return fmt.Errorf("create %q: %w", path, err)
	}

	_, writeErr := file.Write(data)
	closeErr := file.Close()
	if writeErr != nil || closeErr != nil {
		if removeErr := os.Remove(path); removeErr != nil {
			return errors.Join(writeErr, closeErr, fmt.Errorf("remove incomplete file %q: %w", path, removeErr))
		}
		return errors.Join(writeErr, closeErr)
	}
	return nil
}
