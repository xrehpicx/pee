package projectmanager

import (
	"fmt"
	"os"

	projectconfig "github.com/xrehpicx/pee/config"
	"github.com/xrehpicx/pee/controller"
	"github.com/xrehpicx/pee/ui/filepicker"
	"github.com/xrehpicx/pee/ui/table"
	"github.com/xrehpicx/pee/utils"

	btable "github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var ListProjects = &cobra.Command{
	Use:   "ls",
	Short: "List all projects",
	Run: func(cmd *cobra.Command, args []string) {
		projects, err := projectconfig.ListProjects()
		if err != nil {
			fmt.Println(err)
			return
		}
		columns := []btable.Column{
			{Title: "Name", Width: 20},
			{Title: "Session Name", Width: 20},
			{Title: "Working Dir", Width: 50},
			{Title: "Last Opened", Width: 20},
		}
		var rows []btable.Row

		for projectName, config := range projects {
			row := []string{
				projectName,
				config.SessionName,
				config.WorkingDir,
				config.LastOpened.Format("2006-01-02 15:04:05"),
			}
			rows = append(rows, row)
		}
		selectedRow, action := table.Table(columns, rows)
		if action == "edit" {
			// print a vim command to open the config file
			editorCommand, err := projectconfig.GetEditorCommand(selectedRow[0])
			if err != nil {
				editorCommand = "vim"
			}
			utils.EditFile(projectconfig.ProjectConfigFilePath(selectedRow[0]), editorCommand)
			log.Debug("Opened config file", "file", projectconfig.ProjectConfigFilePath(selectedRow[0]))
		}
		if action == "open" {
			ExecuteProjectEnv(selectedRow[0])
		}
	},
}

var RootCmd = &cobra.Command{
	Use:  "pee",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ExecuteProjectEnv(args[0])
	},
}

func ExecuteProjectEnv(projectName string) {
	config, err := projectconfig.GetProjectConfig(projectName)
	if err != nil {
		log.Error(err)
		return
	}
	err = controller.CreateTmuxSession(config)
	if err != nil {
		log.Error(err)
		return
	}
	projectconfig.UpdateLastOpened(projectName)
	log.Debug("Created tmux session", "name", config.SessionName)
}

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
			fmt.Println("Created Project, set up your config by editing:", ppath)
		}
	},
}
