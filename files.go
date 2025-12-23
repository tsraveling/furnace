package main

import (
	"fmt"
	"os"
	"strings"
)

func countLinesInFile(path string) (int, error) {
	// Count existing lines in file
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	lineCount := strings.Count(string(data), "\n")
	if len(data) > 0 && data[len(data)-1] != '\n' {
		lineCount++ // account for last line without newline
	}
	return lineCount, nil
}

func deleteLine(path string, lineNum int) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")

	if lineNum < 1 || lineNum > len(lines) {
		return fmt.Errorf("line %d out of range", lineNum)
	}

	// Remove the line (lineNum is 1-based)
	lines = append(lines[:lineNum-1], lines[lineNum:]...)

	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
}
