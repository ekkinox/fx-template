package cmd

import (
	"fmt"

	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Application version",
	Long:  "Print the application name and version",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := fxconfig.NewDefaultConfigFactory().Create()
		if err != nil {
			fmt.Printf("error getting config: %v", err)
		}

		fmt.Printf("application: %s, version: %s\n", config.AppName(), config.AppVersion())
	},
}
