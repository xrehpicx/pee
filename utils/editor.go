package utils

import (
	"os"
	"os/exec"
)

func EditFile(filePath string, editorCommand string) error {
	if editorCommand == "" {
		editorCommand = "vim"
	}

	cmd := exec.Command(editorCommand, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
