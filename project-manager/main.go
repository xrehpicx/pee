package projectmanager

import (
	projectconfig "github.com/xrehpicx/pee/config"
	"github.com/xrehpicx/pee/controller"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

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
