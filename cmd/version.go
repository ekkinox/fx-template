package cmd

import (
	"fmt"

	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Application version",
	Long:  "Print the application version",
	Run: func(cmd *cobra.Command, args []string) {
		rootFxOpts = append(
			rootFxOpts,
			fx.NopLogger,
			fx.Invoke(func(config *fxconfig.Config) {
				fmt.Printf("version: %s\n", config.AppVersion())
			}),
		)

		app := fx.New(rootFxOpts...)

		if err := app.Start(cmd.Context()); err != nil {
			fmt.Printf("error starting application: %v", err)
		}
		if err := app.Stop(cmd.Context()); err != nil {
			fmt.Printf("error stopping application: %v", err)
		}
	},
}
