package cmd

import (
	"fmt"

	"github.com/CreepySunny/tasker/tasks"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [task description]",
	Short: "Add a new task to your to-do list",
	Long: `Add a new task to your to-do list. Example:

  tasker add "Buy groceries"`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		description := args[0]
		err := tasks.AddTask(fileName, description)
		if err != nil {
			fmt.Printf("Failed to add task: %v\n", err)
			return
		}
		fmt.Println("Task added:", description)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
