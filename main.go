package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var version = "0.0.0"

var targetDir = flag.String("target-dir", ".", "project root directory")

func main() {

	flag.Parse()

	fmt.Printf("github.com/dearing/go-github-releaser version %s\n", version)

	// look for a go-github-releaser.csv file
	file, err := os.Open("go-github-releaser.csv")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}

	defer file.Close()

	// read the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 5 {
			fmt.Printf("Invalid line: %s\n", line)
			continue
		}
		goOS := parts[0]
		goARCH := parts[1]
		name := parts[2]

		do(goOS, goARCH, name)
	}

}

func do(goOS, goARCH, name string) {

	// move to working dir
	origDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	err = os.Chdir(*targetDir)
	if err != nil {
		fmt.Printf("Error changing directory: %v\n", err)
		return
	}

	defer func() {
		err := os.Chdir(origDir)
		if err != nil {
			fmt.Printf("Error restoring original directory: %v\n", err)
		}
	}()

	cmd := exec.Command("go", "build", "-v", "-x", "-o", name)

	// set environment variables
	cmd.Env = append(os.Environ(),
		"GOOS="+goOS,
		"GOARCH="+goARCH,
	)

	log.Printf("Building %s %s/%s\n", name, goOS, goARCH)

	// Save the original stdout and stderr
	origStdout := os.Stdout
	origStderr := os.Stderr

	// Restore stdout and stderr after completion
	defer func() {
		os.Stdout = origStdout
		os.Stderr = origStderr
	}()

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error building: %v\n", err)
		return
	}
	fmt.Printf("Build successful!\n")
}
