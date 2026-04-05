package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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

	triage := fmt.Sprintf(`
# Bug Triage

## Input File
%s

## Raw Bug Report
%s

## Summary
TODO

## Unknowns
- TODO

## Repro Steps
- TODO`, filepath.Base(inputFile), string(bugReportFile))
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
	fmt.Println("triage written:", newFilePath)
}
