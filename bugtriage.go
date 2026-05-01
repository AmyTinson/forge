package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AmyTinson/forge/internal/triage"

	"github.com/google/uuid"
)

/**
TODO:
- declar vars for dir names at top of file for easier changes
- make helper func for file naming
*/

type Run struct {
	ID               string    `json:"id"`
	Workflow         string    `json:"workflow"`
	InputFileName    string    `json:"inputFileName"`
	ArtifactFileName string    `json:"artifactFileName"`
	CreatedAt        time.Time `json:"createdAt"`
	Error            string    `json:"error,omitempty"`
}

func bugTriage(inputFile string) {
	bugReportFile, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("error reading file:", err)
		return
	}

	err = os.MkdirAll("data/artifacts", 0755)
	if err != nil {
		fmt.Println("error making dir:", err)
		return
	}
	err = os.MkdirAll("data/runs", 0755)
	if err != nil {
		fmt.Println("error making dir:", err)
		return
	}
	triageData := triage.TriageData{
		FileName:  filepath.Base(inputFile),
		RawReport: string(bugReportFile),
	}
	triage, err := triage.RenderTriage(triageData)
	if err != nil {
		fmt.Println("renderTriage:", err)
		return
	}
	inputFileName := filepath.Base(inputFile)
	inputFileExt := filepath.Ext(inputFileName)
	newFileName := strings.TrimSuffix(inputFileName, inputFileExt) + "-triage.md"
	newFilePath := "data/artifacts/" + newFileName

	artifactFile, err := os.OpenFile(newFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("error opening file:", err)
		return
	}
	defer artifactFile.Close()

	_, err = artifactFile.WriteString(triage + "\n")
	if err != nil {
		fmt.Println("error writing to file:", err)
		return
	}
	run := Run{
		ID:               uuid.New().String(),
		Workflow:         "bug-triage",
		InputFileName:    inputFileName,
		ArtifactFileName: newFileName,
		CreatedAt:        time.Now(),
	}
	runFilePath := "data/runs/" + run.ID + ".json"
	runFile, err := os.OpenFile(runFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("error opening run file:", err)
		return
	}
	defer runFile.Close()

	runData, err := json.Marshal(run)
	if err != nil {
		fmt.Println("error marshaling run data:", err)
		return
	}

	_, err = runFile.Write(runData)
	if err != nil {
		fmt.Println("error writing run data to file:", err)
		return
	}
	fmt.Println("triage written:", newFilePath)
}
