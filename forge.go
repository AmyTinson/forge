package main

import (
	"log"
	"os"
)

func forge() {
	args := os.Args
	if len(args) < 3 {
		log.Fatal("usage: forge <command> <arg>")
	}
	command := os.Args[1]
	arg := os.Args[2]

	switch command {
	case "read":
		{
			file, err := os.ReadFile(arg)
			if err != nil {
				log.Fatal(err)
			}
			os.Stdout.Write(file)
		}

	case "run":
		{
			switch arg {
			case "bug-triage":
				{
					if len(args) < 5 || args[3] != "--input" || args[4] == "" {
						log.Fatal("usage: forge run bug-triage --input <file>")
					}
					bugTriage(args[4])
				}
			default:
				log.Fatal("Unknown argument for run command: " + arg)
			}
		}

	case "note":
		{
			switch arg {
			case "add":
				{
					if len(args) < 4 || args[3] == "" {
						log.Fatal("usage: forge note add \"your note\"")
					}
					addNote(args[3])
				}
			default:
				log.Fatal("Unknown argument for add command: " + arg)
			}
		}
	default:
		log.Fatal("Unknown command: " + command)
	}

}
