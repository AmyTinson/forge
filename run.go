package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Run struct {
	ID               string    `json:"id"`
	Workflow         string    `json:"workflow"`
	InputFileName    string    `json:"inputFileName"`
	ArtifactFileName string    `json:"artifactFileName"`
	CreatedAt        time.Time `json:"createdAt"`
	Error            string    `json:"error,omitempty"`
}

func persistRun(run Run, artifact []byte) (string, error) {
	artifactDir := filepath.Join("data", "artifacts")
	runDir := filepath.Join("data", "runs")
	if err := os.MkdirAll(artifactDir, 0755); err != nil {
		return "", fmt.Errorf("create artifact directory: %w", err)
	}
	if err := os.MkdirAll(runDir, 0755); err != nil {
		return "", fmt.Errorf("create run directory: %w", err)
	}

	runData, err := json.Marshal(run)
	if err != nil {
		return "", fmt.Errorf("marshal run record: %w", err)
	}

	artifactPath := filepath.Join(artifactDir, run.ArtifactFileName)
	if err := writeNewFile(artifactPath, artifact); err != nil {
		return "", fmt.Errorf("write artifact: %w", err)
	}

	runPath := filepath.Join(runDir, run.ID+".json")
	if err := writeNewFile(runPath, runData); err != nil {
		if cleanupErr := os.Remove(artifactPath); cleanupErr != nil {
			return "", errors.Join(
				fmt.Errorf("write run record: %w", err),
				fmt.Errorf("remove incomplete artifact %q: %w", artifactPath, cleanupErr),
			)
		}
		return "", fmt.Errorf("write run record: %w", err)
	}
	return artifactPath, nil
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
