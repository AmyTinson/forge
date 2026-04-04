package main

import (
	"fmt"
	"os"
)

func addNote(note string) {
	f, err := os.OpenFile("notes/notes.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer f.Close()

	_, err = f.WriteString(note + "\n")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("note added:", note)
}
