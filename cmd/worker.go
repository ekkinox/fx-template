package cmd

import (
	"github.com/ekkinox/fx-template/internal/worker/pubsub"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/ekkinox/fx-template/modules/fxpubsub"
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

		worker := fx.New(
			// core
			fxconfig.FxConfigModule,
			fxlogger.FxLoggerModule,
			fxtracer.FxTracerModule,
			fxhealthchecker.FxHealthCheckerModule,
			// worker
			fxpubsub.FxPubSubModule,
			fx.Provide(pubsub.NewSubscribeWorker),
			fx.Invoke(func(w *pubsub.SubscribeWorker) {
				w.Run(cmd.Context())
			}),
			// logger
			fx.WithLogger(fxlogger.FxEventLogger),
		)

		worker.Run()
	},
}
