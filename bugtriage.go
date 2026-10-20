package main

import (
	"context"
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

type BugTriageWorkflow struct {
	InputFile string
}

var _ Workflow = BugTriageWorkflow{}

func (w BugTriageWorkflow) Name() string {
	return bugTriageWorkflowName
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
	artifactPath, err := persistRun(run, []byte(renderedTriage+"\n"))
	if err != nil {
		return fmt.Errorf("persist triage run: %w", err)
	}

	fmt.Println("triage written:", artifactPath)
	return nil
}
