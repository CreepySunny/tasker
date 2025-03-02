/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tasker",
	Short: "Convenient To-Do List CLI",
	Long: `Tasker is a command-line interface (CLI) tool designed to help you manage your to-do lists efficiently.
With Tasker, you can quickly add, remove, and view tasks directly from your terminal, making it easier to stay organized and productive.

Features:
- Add new tasks with a simple command
- Remove completed or unnecessary tasks
- View all your tasks in a clear and concise format
- Mark tasks as completed

Examples:
- Add a new task: tasker add "Buy groceries"
- Remove a task: tasker remove 1
- View all tasks: tasker list
- Mark a task as completed: tasker complete 1`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tasker.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
