package cmd

import (
	"github.com/ekkinox/fx-template/internal/app"
	"github.com/ekkinox/fx-template/modules/fxtracer"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve application",
	Long:  "Serve the application server",
	Run: func(cmd *cobra.Command, args []string) {
		rootFxOpts = append(
			rootFxOpts,
			fxtracer.FxTracerModule,
			app.RegisterModules(),
			app.RegisterHandlers(),
			app.RegisterServices(),
			app.RegisterOverrides(),
		)

		fx.New(rootFxOpts...).Run()
	},
}
