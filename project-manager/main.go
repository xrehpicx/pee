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

var KillTmuxSessionCmd = &cobra.Command{
	Use:  "kill",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sessionName := args[0]
		err := controller.KillTmuxSession(sessionName)
		if err != nil {
			log.Error(err)
			return
		}
		log.Debug("Killed tmux session", "name", sessionName)
	},
}

func init() {
	RootCmd.AddCommand(KillTmuxSessionCmd)
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
