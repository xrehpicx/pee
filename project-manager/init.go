package projectmanager

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	projectconfig "github.com/xrehpicx/pee/config"
	"github.com/xrehpicx/pee/ui/filepicker"
	"github.com/xrehpicx/pee/utils"
)

var InitCmd = &cobra.Command{
	Use:  "init",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		projectExists := projectconfig.ProjectExists(projectName)
		if projectExists {
			log.Warn("Project already exists", "name", projectName)
			return
		}

		selected, err := filepicker.FilePicker("Select your project dir", "Selected Dir: ")
		if selected == "" {
			log.Warn("No dir selected, aborting")
			return
		}
		if err != nil {
			log.Error(err)
			return
		}
		log.Debug("Selected", "work_dir", selected)

		// Define the session configuration
		workingDir := selected
		windows := []projectconfig.Window{
			{
				WindowName: "editor",
				Layout:     "8070,202x58,0,0[202x46,0,0,89,202x11,0,47,92]",
				Panes: []projectconfig.Pane{
					{
						ShellCommand: []string{"echo 'command to open your editor'"},
					},
					{
						ShellCommand: []string{"echo 'run dev server'"},
					},
				},
			},
			{
				WindowName: "ssh windows",
				ShellCommandBefore: []string{
					"ls -lr",
				},
				Panes: []projectconfig.Pane{
					{
						ShellCommand: []string{"echo 'command to open your ssh windows'"},
					},
					{
						ShellCommand: []string{"echo 'command to open your ssh windows'"},
					},
				},
			},
			{
				WindowName: "git",
				Panes: []projectconfig.Pane{
					{
						ShellCommand: []string{"echo 'command to open your git client'"},
					},
				},
			},
		}

		logger := log.NewWithOptions(os.Stderr, log.Options{
			ReportCaller:    false,
			ReportTimestamp: false,
		})

		ppath, err := projectconfig.CreateProject(projectName, workingDir, windows)
		if err != nil {
			logger.Error(err)
		} else {
			editorCommand, err := projectconfig.GetEditorCommand(projectName)
			if err != nil {
				editorCommand = ""
			}
			utils.EditFile(ppath, editorCommand)
			fmt.Println("Created Project, config is at:", ppath)
		}
	},
}
