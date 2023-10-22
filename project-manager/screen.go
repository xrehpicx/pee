package projectmanager

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	projectconfig "github.com/xrehpicx/pee/config"
	"github.com/xrehpicx/pee/controller"
)

var KillScreenCmd = &cobra.Command{
	Use:  "kill",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		controller.KillScreenSession(args[0])
	},
}

var ScreenCmd = &cobra.Command{
	Use:  "scr",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ExecuteProjectEnvUsingScreen(args[0])
	},
}

func ExecuteProjectEnvUsingScreen(projectName string) {
	config, err := projectconfig.GetProjectConfig(projectName)
	if err != nil {
		log.Error(err)
		return
	}
	err = controller.CreateScreenSession(config)
	if err != nil {
		log.Error(err)
		return
	}
	projectconfig.UpdateLastOpened(projectName)
	log.Debug("Created tmux session", "name", config.SessionName)
}

func init() {
	ScreenCmd.AddCommand(KillScreenCmd)
}
