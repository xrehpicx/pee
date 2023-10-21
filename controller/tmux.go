package controller

import (
	"fmt"
	"os/exec"
	projectconfig "pee/config"
	"strings"

	"github.com/charmbracelet/log"
)

// CreateTmuxSession creates a tmux session based on the given Configuration.
func CreateTmuxSession(config *projectconfig.Configuration) error {
	sessionName := config.SessionName

	// Check if the session exists
	checkSessionCmd := exec.Command("tmux", "has-session", "-t", sessionName)
	if err := checkSessionCmd.Run(); err == nil {
		// If it exists, switch to the session
		switchSessionCmd := exec.Command("tmux", "switch-client", "-t", sessionName)
		if err := switchSessionCmd.Run(); err != nil {
			return err
		}
	} else {
		// If it doesn't exist, create the session
		createSessionCmd := exec.Command("tmux", "new-session", "-d", "-s", sessionName)
		if err := createSessionCmd.Run(); err != nil {
			return err
		}
		log.Info("Ran command", "command", createSessionCmd.String())

		// Change the working directory
		changeDirCmd := exec.Command("tmux", "send-keys", "-t", sessionName, "cd "+config.WorkingDir, "Enter")
		if err := changeDirCmd.Run(); err != nil {
			return err
		}
		log.Info("Ran command", "command", changeDirCmd.String())

		// Send commands to the session for the first tab
		sendCommandsCmd := exec.Command("tmux", "send-keys", "-t", sessionName, strings.Join(config.Tabs[0].Commands, " && "), "Enter")
		if err := sendCommandsCmd.Run(); err != nil {
			return err
		}
		log.Info("Ran command", "command", sendCommandsCmd.String())
		// Rename the tab to the specified name
		renameTabCmd := exec.Command("tmux", "rename-window", "-t", sessionName+":1", config.Tabs[0].Name)
		if err := renameTabCmd.Run(); err != nil {
			return err
		}

		// Create and run commands for additional tabs
		for i, tab := range config.Tabs[1:] {
			windowName := fmt.Sprintf("%s:%d", sessionName, i+2)
			createWindowCmd := exec.Command("tmux", "new-window", "-t", windowName, "-n", tab.Name)
			if err := createWindowCmd.Run(); err != nil {
				return err
			}
			log.Info("Ran command", "command", createWindowCmd.String())

			changeDirCmd = exec.Command("tmux", "send-keys", "-t", windowName, "cd "+config.WorkingDir, "Enter")
			if err := changeDirCmd.Run(); err != nil {
				return err
			}
			log.Info("Ran command", "command", changeDirCmd.String())

			sendCommandsCmd = exec.Command("tmux", "send-keys", "-t", windowName, strings.Join(tab.Commands, " && "), "Enter")
			if err := sendCommandsCmd.Run(); err != nil {
				return err
			}
			log.Info("Ran command", "command", sendCommandsCmd.String())
		}

		// Select the initial window and switch to the session
		selectWindowCmd := exec.Command("tmux", "select-window", "-t", sessionName+":1")
		if err := selectWindowCmd.Run(); err != nil {
			return err
		}

		switchSessionCmd := exec.Command("tmux", "switch-client", "-t", sessionName)
		if err := switchSessionCmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
