package main

import (
	"fmt"
	"os"

	projectconfig "github.com/xrehpicx/pee/config"
	projectmanager "github.com/xrehpicx/pee/project-manager"
)

func init() {
	projectconfig.Init()
	projectmanager.RootCmd.AddCommand(projectmanager.ListProjects)
	projectmanager.RootCmd.AddCommand(projectmanager.InitCmd)
	projectmanager.RootCmd.AddCommand(projectmanager.ScreenCmd)
}

func main() {
	if err := projectmanager.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
