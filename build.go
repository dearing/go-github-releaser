package main

import (
	"os"
	"os/exec"
)

func doBuild(goOS, goARCH, target string) error {

	cmd := exec.Command("go", "build", "-o", target, *srcDir)

	// set environment variables
	cmd.Env = append(os.Environ(),
		"GOOS="+goOS,
		"GOARCH="+goARCH,
	)

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
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
