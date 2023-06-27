package cmd

import (
	"github.com/ekkinox/fx-template/internal/grpc"
	"github.com/ekkinox/fx-template/internal/server"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxgrpcserver"
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

		server := fx.New(
			// core
			fxconfig.FxConfigModule,
			fxlogger.FxLoggerModule,
			fxtracer.FxTracerModule,
			// http
			fxhttpserver.FxHttpServerModule,
			server.RegisterModules(),
			server.RegisterHandlers(),
			server.RegisterServices(),
			server.RegisterOverrides(),
			// grpc
			fxgrpcserver.FxGrpcServerModule,
			grpc.RegisterGrpcServices(),
			// logger
			fx.WithLogger(fxlogger.FxEventLogger),
		)

		server.Run()
	},
}
