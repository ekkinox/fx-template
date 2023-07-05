package cmd

import (
	"github.com/ekkinox/fx-template/internal/server"
	"github.com/ekkinox/fx-template/internal/server/grpc"
	"github.com/ekkinox/fx-template/internal/server/http"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxgrpcserver"
	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/ekkinox/fx-template/modules/fxtracer"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Server application",
	Long:  "HTTP Server application",
	Run: func(cmd *cobra.Command, args []string) {

		ctx := cmd.Context()

		server := fx.New(
			// logger
			fx.WithLogger(fxlogger.FxEventLogger),
			// core
			fxconfig.FxConfigModule,
			fxlogger.FxLoggerModule,
			fxtracer.FxTracerModule,
			fxhealthchecker.FxHealthCheckerModule,
			// common
			server.RegisterModules(ctx),
			server.RegisterServices(ctx),
			server.RegisterOverrides(ctx),
			// http server
			fxhttpserver.FxHttpServerModule,
			http.RegisterHandlers(),
			// grpc server
			fxgrpcserver.FxGrpcServerModule,
			grpc.RegisterGrpcServices(),
		)

		server.Run()
	},
}
