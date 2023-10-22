package projectmanager

import (
	"fmt"

	btable "github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	projectconfig "github.com/xrehpicx/pee/config"
	"github.com/xrehpicx/pee/ui/table"
	"github.com/xrehpicx/pee/utils"
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
				editorCommand = ""
			}
			utils.EditFile(projectconfig.ProjectConfigFilePath(selectedRow[0]), editorCommand)
			log.Debug("Opened config file", "file", projectconfig.ProjectConfigFilePath(selectedRow[0]))
		}
		if action == "open" {
			ExecuteProjectEnv(selectedRow[0])
		}
	},
}
