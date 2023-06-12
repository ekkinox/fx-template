package cmd

import (
	"fmt"
	"os"

	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var rootFxOpts []fx.Option

var rootCmd = &cobra.Command{
	Use: "app",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		rootFxOpts = append(
			rootFxOpts,
			fxconfig.FxConfigModule,
			fxlogger.FxLoggerModule,
			fx.WithLogger(fxlogger.FxEventLogger),
		)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
