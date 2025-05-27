package cmd

import (
	"fmt"

	"github.com/CreepySunny/tasker/tasks"
	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "complete [task ID]",
	Short: "Mark a task as completed",
	Long: `Mark a task as completed. Example:
	  tasker complete 1`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskID := args[0]
		err := tasks.CompleteTask(fileName, taskID)
		if err != nil {
			fmt.Printf("Failed to complete task: %v\n", err)
			return
		}
		fmt.Println("TaskID completd:", taskID)
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
