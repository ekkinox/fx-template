package cmd

import (
	"github.com/ekkinox/fx-template/internal/worker"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/ekkinox/fx-template/modules/fxtracer"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func init() {
	rootCmd.AddCommand(workerCmd)
}

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Worker application",
	Long:  "Pub/Sub worker application",
	Run: func(cmd *cobra.Command, args []string) {

		ctx := cmd.Context()

		worker := fx.New(
			// logger
			fx.WithLogger(fxlogger.FxEventLogger),
			// core
			fxconfig.FxConfigModule,
			fxlogger.FxLoggerModule,
			fxtracer.FxTracerModule,
			fxhealthchecker.FxHealthCheckerModule,
			// worker
			worker.RegisterModules(ctx),
			worker.RegisterServices(ctx),
			worker.RegisterOverrides(ctx),
		)

		worker.Run()
	},
}
