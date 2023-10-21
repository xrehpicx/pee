package projectmanager

import (
	"fmt"
	"os"

	projectconfig "github.com/xrehpicx/pee/config"
	"github.com/xrehpicx/pee/controller"
	"github.com/xrehpicx/pee/ui/filepicker"
	"github.com/xrehpicx/pee/ui/table"

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
			fmt.Println("vim", projectconfig.ProjectConfigFilePath(selectedRow[0]))
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
	log.Info("Created tmux session", "name", config.SessionName)
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
		log.Info("Selected", "work_dir", selected)

		sessionName := projectName
		tabs := []struct {
			Name     string
			Commands []string
		}{
			{
				Name:     "editor",
				Commands: []string{"echo 'command to open ur editor'"},
			},
			{
				Name:     "dev server",
				Commands: []string{"echo 'command to start dev server'", "echo 'command to just initialize ur dependencies'"},
			},
			{
				Name:     "git",
				Commands: []string{"echo 'command to open ur git client (use lazygit its amazing)'"},
			},
		}
		logger := log.NewWithOptions(os.Stderr, log.Options{
			ReportCaller:    false,
			ReportTimestamp: false,
		})
		ppath, err := projectconfig.CreateProject(projectName, sessionName, selected, tabs)
		if err != nil {
			logger.Error(err)
		} else {
			// logger.Info("Created Project", "path", ppath)
			fmt.Println("Created Project", "setup your config by editing: ", ppath)
		}
	},
}
