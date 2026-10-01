package main

import (
	"context"
	"fmt"
	"os"
)

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
	ctx := context.Background()
	if len(args) < 2 {
		return fmt.Errorf("usage: forge <command> <arg>")
	}
	command := args[0]
	arg := args[1]

	switch command {
	case "read":
		return runRead(arg)

	case "run":

		switch arg {
		case "bug-triage":
			if len(args) < 4 || args[2] != "--input" || args[3] == "" {
				return fmt.Errorf("usage: forge run bug-triage --input <file>")
			}
			workflow := BugTriageWorkflow{
				InputFile: args[3],
			}

			err := workflow.Run(ctx)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown argument for run command: %s", arg)
		}

	case "note":

		switch arg {
		case "add":

			if len(args) < 3 || args[2] == "" {
				return fmt.Errorf("usage: forge note add \"your note\"")
			}
			addNote(args[2])

		default:
			return fmt.Errorf("unknown argument for add command: %s", arg)
		}

	default:
		return fmt.Errorf("unknown command: %s", command)
	}
	return nil

}
