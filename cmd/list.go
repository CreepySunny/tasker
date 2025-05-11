package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/CreepySunny/tasker/tasks"
	"github.com/spf13/cobra"
)

var showAll bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Long: `List all tasks in your task manager. 
You can use the --all or -a flag to include completed tasks in the list. 
For example:
  tasker list --all
This will show both completed and pending tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := tasks.ListTasks(fileName, showAll)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: %w\n", err)
			return
		}

		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

		for _, task := range tasks {
			fmt.Fprintln(tw, task.String())
		}

		tw.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all tasks")
}
