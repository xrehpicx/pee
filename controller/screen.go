package controller

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
	projectconfig "github.com/xrehpicx/pee/config"
)

// RunShellCommand executes a shell command and logs it.
func RunShellCommand(cmd *exec.Cmd) error {
	log.Debug("Running command:", "command", cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("Error running command:", "command", cmd.String(), "error", err)
		log.Info("Command output:", "output", string(output))
		return err
	}
	return nil
}

// CreateScreenSession creates a new screen session based on the given Configuration.
func CreateScreenSession(config *projectconfig.Configuration) error {
	sessionName := config.SessionName

	// Check if the session already exists
	checkSessionCmd := exec.Command("screen", "-S", sessionName, "-Q", "windows")
	_, err := checkSessionCmd.CombinedOutput()
	if err == nil {
		// If it exists, attach to the session
		attachSessionCmd := exec.Command("screen", "-d", "-r", sessionName)
		RunShellCommand(attachSessionCmd)
	} else {
		// If it doesn't exist, create the session
		createSessionCmd := exec.Command("screen", "-S", sessionName, "-d", "-m")
		RunShellCommand(createSessionCmd)

		// Create and run commands for windows
		for i, window := range config.Windows {
			windowName := fmt.Sprintf("%s-%d", sessionName, i+1)

			// Create a new window within the session
			createWindowCmd := exec.Command("screen", "-S", sessionName, "-X", "screen", "-t", windowName)
			RunShellCommand(createWindowCmd)

			// Change the working directory for the window
			changeDirCmd := exec.Command("screen", "-S", sessionName, "-p", fmt.Sprint(i), "-X", "chdir", config.WorkingDir)
			RunShellCommand(changeDirCmd)

			// Send commands to the window
			sendCommandsCmd := exec.Command("screen", "-S", sessionName, "-p", fmt.Sprint(i), "-X", "stuff", strings.Join(window.Panes[0].ShellCommand, " && ")+"\n")
			RunShellCommand(sendCommandsCmd)

			// Rename the window to the specified name
			renameWindowCmd := exec.Command("screen", "-S", sessionName, "-p", fmt.Sprint(i), "-X", "title", window.WindowName)
			RunShellCommand(renameWindowCmd)

			// warn user of compatibility issues using more than one pane with screen
			if len(window.Panes) > 1 {
				log.Warn("Screen does not support multiple panes. Only the first pane will be used.")
			}

			// Create and run commands for additional panes in the window
			for _, pane := range window.Panes[1:] {
				// Split the window vertically
				splitPaneCmd := exec.Command("screen", "-S", sessionName, "-p", fmt.Sprint(i), "-X", "split")
				RunShellCommand(splitPaneCmd)

				// Select the new pane
				selectPaneCmd := exec.Command("screen", "-S", sessionName, "-p", fmt.Sprint(i+1)) // Select the next pane
				RunShellCommand(selectPaneCmd)

				// Send commands to the new pane
				sendCommandsCmd := exec.Command("screen", "-S", sessionName, "-X", "stuff", strings.Join(pane.ShellCommand, " && ")+"\n")
				RunShellCommand(sendCommandsCmd)
			}

			if window.Layout != "" {
				layoutCmd := exec.Command("screen", "-S", sessionName, "-p", fmt.Sprint(i), "-X", "layout", window.Layout)
				RunShellCommand(layoutCmd)
			}
		}

		if config.Attach {
			// Attach to the session
			attachSessionCmd := exec.Command("screen", "-d", "-r", sessionName)
			RunShellCommand(attachSessionCmd)
		}
	}

	return nil
}

// KillScreenSession kills a screen session by name.
func KillScreenSession(sessionName string) error {
	// Kill the screen session
	killSessionCmd := exec.Command("screen", "-S", sessionName, "-X", "quit")
	err := killSessionCmd.Run()
	if err != nil {
		log.Error("Error killing screen session:", "sessionName", sessionName, "error", err)
		return err
	}
	log.Info("Killed screen session:", "sessionName", sessionName)
	return nil
}
