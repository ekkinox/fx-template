package cmd

import (
	"github.com/ekkinox/fx-template/internal/server"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Server application",
	Long:  "HTTP Server application",
	Run: func(cmd *cobra.Command, args []string) {
		server.BootstrapServer(cmd.Context()).Run()
	},
}
