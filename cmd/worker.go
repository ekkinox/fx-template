package cmd

import (
	"github.com/ekkinox/fx-template/internal/worker"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(workerCmd)
}

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Worker application",
	Long:  "Pub/Sub worker application",
	Run: func(cmd *cobra.Command, args []string) {
		worker.BootstrapWorker(cmd.Context()).Run()
	},
}
