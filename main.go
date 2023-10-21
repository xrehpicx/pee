package main

import (
	"fmt"
	"os"
	projectconfig "pee/config"
	projectmanager "pee/project-manager"
)

func init() {
	projectconfig.Init()
	projectmanager.RootCmd.AddCommand(projectmanager.ListProjects)
	projectmanager.RootCmd.AddCommand(projectmanager.InitCmd)
}

func main() {
	if err := projectmanager.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
