package main

import (
	"fmt"
	"os"
)

func addNote(note string) error {
	f, err := os.OpenFile("notes/notes.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open notes file: %w", err)
	}

	_, err = f.WriteString(note + "\n")
	closeErr := f.Close()
	if err != nil {
		return fmt.Errorf("write notes file: %w", err)
	}
	if closeErr != nil {
		return fmt.Errorf("close notes file: %w", closeErr)
	}
	return nil
}
