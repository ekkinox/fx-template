package cmd

import (
	"github.com/ekkinox/fx-template/internal/worker/pubsub"
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
		rootFxOpts = append(
			rootFxOpts,
			fxtracer.FxTracerModule,
			fxpubsub.FxPubSubModule,
			fx.Provide(pubsub.NewSubscribeWorker),
			fx.Invoke(func(worker *pubsub.SubscribeWorker) {
				worker.Run()
			}),
		)

		fx.New(rootFxOpts...).Run()
	},
}
