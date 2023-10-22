package utils

import (
	"os"
	"os/exec"
)

func checkNvimExists() bool {
	cmd := exec.Command("nvim", "--version")
	err := cmd.Run()
	return err == nil
}

func EditFile(filePath string, editorCommand string) error {
	if editorCommand == "" {
		if checkNvimExists() {
			editorCommand = "nvim"
		} else {
			editorCommand = "vim"
		}
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
