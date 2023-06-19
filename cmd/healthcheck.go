package cmd

import (
	"fmt"
	"os"

	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxgorm"
	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/ekkinox/fx-template/modules/fxpubsub"
	"github.com/ekkinox/fx-template/modules/fxtracer"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func init() {
	rootCmd.AddCommand(healthCheckCmd)
}

var healthCheckCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "health check",
	Long:  "health check for application",
	Run: func(cmd *cobra.Command, args []string) {

		success := true

		healthcheck := fx.New(
			fxconfig.FxConfigModule,
			fxlogger.FxLoggerModule,
			fxtracer.FxTracerModule,
			fxhealthchecker.FxHealthCheckerModule,
			fxgorm.FxGormModule,
			fxpubsub.FxPubSubModule,
			fx.Provide(
				fxhealthchecker.AsProbe(fxgorm.NewGormProbe),
				fxhealthchecker.AsProbe(fxpubsub.NewPubSubProbe),
			),
			fx.Invoke(func(c *fxhealthchecker.Checker, l *fxlogger.Logger) {
				result := c.Run(cmd.Context())

				for n, r := range result.ProbesResults {
					fmt.Printf("[%s] success: %v, message: %sn", n, r.Success, r.Message)
				}
			}),
			fx.NopLogger,
		)

		if err := healthcheck.Start(cmd.Context()); err != nil {
			fmt.Printf("error starting health check: %v", err)
			os.Exit(1)
		}

		if err := healthcheck.Stop(cmd.Context()); err != nil {
			fmt.Printf("error stopping health check: %v", err)
			os.Exit(1)
		}

		if !success {
			fmt.Println("== health check failed ==")
			os.Exit(1)
		}

		fmt.Println("== health check success ==")
		os.Exit(0)
	},
}
