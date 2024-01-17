package controller

import (
	"fmt"
	"os/exec"
	"strings"

	projectconfig "github.com/xrehpicx/pee/config"

	"github.com/charmbracelet/log"
)

// CreateTmuxSession creates a Tmux session based on the given Configuration.
func CreateTmuxSession(config *projectconfig.Configuration) error {
	sessionName := config.SessionName

	// Check if the session exists
	checkSessionCmd := exec.Command("tmux", "has-session", "-t", sessionName)
	log.Debug("Ran command", "command", checkSessionCmd.String())
	if err := checkSessionCmd.Run(); err == nil {
		// If it exists, switch to the session
		switchSessionCmd := exec.Command("tmux", "switch-client", "-t", sessionName)
		if err := switchSessionCmd.Run(); err != nil {
			return err
		}
		log.Debug("Ran command", "command", switchSessionCmd.String())
	} else {
		// If it doesn't exist, create the session
		createSessionCmd := exec.Command("tmux", "new-session", "-d", "-s", sessionName, "-c", config.WorkingDir)
		// createSessionCmd := exec.Command("tmux", "new-session", "-d", "-s", sessionName)
		if err := createSessionCmd.Run(); err != nil {
			return err
		}
		log.Debug("Ran command", "command", createSessionCmd.String())

		// Change the working directory
		changeDirCmd := exec.Command("tmux", "send-keys", "-t", sessionName, "cd "+config.WorkingDir, "Enter")
		if err := changeDirCmd.Run(); err != nil {
			return err
		}
		log.Debug("Ran command", "command", changeDirCmd.String())

		// Create the first window outside the loop
		// createWindow(config, sessionName, 0)

		window := config.Windows[0]
		windowName := fmt.Sprintf("%s:%d", sessionName, 1)

		sendCommandsCmd := exec.Command("tmux", "send-keys", "-t", windowName, strings.Join(window.Panes[0].ShellCommand, " && "), "Enter")
		if err := sendCommandsCmd.Run(); err != nil {
			return err
		}
		log.Debug("Ran command", "command", sendCommandsCmd.String())

		// Rename the window to the specified name
		renameWindowCmd := exec.Command("tmux", "rename-window", "-t", windowName, window.WindowName)
		if err := renameWindowCmd.Run(); err != nil {
			return err
		}
		log.Debug("Ran command", "command", renameWindowCmd.String())

		// Create and run commands for additional panes in the window
		for j, pane := range window.Panes[1:] {
			paneName := fmt.Sprintf("%s:%d.%d", sessionName, 1, j+2)

			// Split the window horizontally
			splitPaneCmd := exec.Command("tmux", "split-window", "-t", windowName, "-h", "-p", "50")
			if err := splitPaneCmd.Run(); err != nil {
				return err
			}
			log.Debug("Ran command", "command", splitPaneCmd.String())

			// Select the new pane
			selectPaneCmd := exec.Command("tmux", "select-pane", "-t", paneName)
			if err := selectPaneCmd.Run(); err != nil {
				return err
			}
			log.Debug("Ran command", "command", selectPaneCmd.String())

			// Change the working directory
			changeDirCmd := exec.Command("tmux", "send-keys", "-t", paneName, "cd "+config.WorkingDir, "Enter")
			if err := changeDirCmd.Run(); err != nil {
				return err
			}
			log.Debug("Ran command", "command", changeDirCmd.String())

			// Send commands to the pane
			sendCommandsCmd := exec.Command("tmux", "send-keys", "-t", paneName, strings.Join(pane.ShellCommand, " && "), "Enter")
			if err := sendCommandsCmd.Run(); err != nil {
				return err
			}
			log.Debug("Ran command", "command", sendCommandsCmd.String())
		}

		if window.Layout != "" {
			layoutCmd := exec.Command("tmux", "select-layout", "-t", windowName, window.Layout)
			if err := layoutCmd.Run(); err != nil {
				return err
			}
			log.Debug("Ran command", "command", layoutCmd.String())
		}

		// Create and run commands for each window inside the loop
		for i := 1; i < len(config.Windows); i++ {
			createWindow(config, sessionName, i)
		}

		// Select the initial window and switch to the session
		defaultWindow := sessionName + ":1"
		if config.StartupWindow != "" {
			defaultWindow = config.StartupWindow
		}

		selectWindowCmd := exec.Command("tmux", "select-window", "-t", defaultWindow)
		if err := selectWindowCmd.Run(); err != nil {
			return err
		}
		log.Debug("Ran command", "command", selectWindowCmd.String())

		// Select initial pane
		if config.StartupPane > 0 {
			defaultPane := fmt.Sprintf("%s:%d.%d", sessionName, config.StartupPane, 1)
			selectPaneCmd := exec.Command("tmux", "select-pane", "-t", defaultPane)
			if err := selectPaneCmd.Run(); err != nil {
				return err
			}
			log.Debug("Ran command", "command", selectPaneCmd.String())
		}

		if config.Attach {
			switchSessionCmd := exec.Command("tmux", "switch-client", "-t", sessionName)
			if err := switchSessionCmd.Run(); err != nil {
				return err
			}
			log.Debug("Ran command", "command", switchSessionCmd.String())
		}
	}

	return nil
}

func createWindow(config *projectconfig.Configuration, sessionName string, index int) error {
	if index >= len(config.Windows) {
		return nil
	}
	window := config.Windows[index]
	windowName := fmt.Sprintf("%s:%d", sessionName, index+1)

	// Create a new window
	createWindowCmd := exec.Command("tmux", "new-window", "-t", sessionName, "-n", windowName)
	if err := createWindowCmd.Run(); err != nil {
		return err
	}
	log.Debug("Ran command", "command", createWindowCmd.String())

	// Change the working directory for the window
	changeDirCmd := exec.Command("tmux", "send-keys", "-t", windowName, "cd "+config.WorkingDir, "Enter")
	if err := changeDirCmd.Run(); err != nil {
		return err
	}
	log.Debug("Ran command", "command", changeDirCmd.String())

	// Send commands to the window
	sendCommandsCmd := exec.Command("tmux", "send-keys", "-t", windowName, strings.Join(window.Panes[0].ShellCommand, " && "), "Enter")
	if err := sendCommandsCmd.Run(); err != nil {
		return err
	}
	log.Debug("Ran command", "command", sendCommandsCmd.String())

	// Rename the window to the specified name
	renameWindowCmd := exec.Command("tmux", "rename-window", "-t", windowName, window.WindowName)
	if err := renameWindowCmd.Run(); err != nil {
		return err
	}
	log.Debug("Ran command", "command", renameWindowCmd.String())

	// Create and run commands for additional panes in the window
	for j, pane := range window.Panes[1:] {
		paneName := fmt.Sprintf("%s:%d.%d", sessionName, index+1, j+2)

		// Split the window horizontally
		splitPaneCmd := exec.Command("tmux", "split-window", "-t", windowName, "-h", "-p", "50")
		if err := splitPaneCmd.Run(); err != nil {
			return err
		}
		log.Debug("Ran command", "command", splitPaneCmd.String())

		// Select the new pane
		selectPaneCmd := exec.Command("tmux", "select-pane", "-t", paneName)
		if err := selectPaneCmd.Run(); err != nil {
			return err
		}
		log.Debug("Ran command", "command", selectPaneCmd.String())

		// Change the working directory for the pane
		changeDirCmd := exec.Command("tmux", "send-keys", "-t", paneName, "cd "+config.WorkingDir, "Enter")
		if err := changeDirCmd.Run(); err != nil {
			return err
		}
		log.Debug("Ran command", "command", changeDirCmd.String())

		// Send commands to the pane
		sendCommandsCmd := exec.Command("tmux", "send-keys", "-t", paneName, strings.Join(pane.ShellCommand, " && "), "Enter")
		if err := sendCommandsCmd.Run(); err != nil {
			return err
		}
		log.Debug("Ran command", "command", sendCommandsCmd.String())
	}

	if window.Layout != "" {

		layoutCmd := exec.Command("tmux", "select-layout", "-t", windowName, window.Layout)
		if err := layoutCmd.Run(); err != nil {
			return err
		}
		log.Debug("Ran command", "command", layoutCmd.String())
	}

	return nil
}

func KillTmuxSession(sessionName string) error {
	killSessionCmd := exec.Command("tmux", "kill-session", "-t", sessionName)
	if err := killSessionCmd.Run(); err != nil {
		return err
	}
	log.Debug("Ran command", "command", killSessionCmd.String())
	return nil
}
