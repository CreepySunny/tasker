package cmd

import (
	"fmt"

	"github.com/CreepySunny/tasker/tasks"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [task ID]",
	Short: "Delete a task by its ID",
	Long: `Delete a task from your to-do list by its ID. Example:
  tasker delete 1`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskID := args[0]
		err := tasks.DeleteTask(fileName, taskID)
		if err != nil {
			fmt.Printf("Failed to delete task: %v\n", err)
			return
		}
		fmt.Println("Task deleted:", taskID)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
