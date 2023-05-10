package services

import (
	"github.com/labstack/echo/v4"
	"strings"

	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxlogger"
)

type TestService struct {
	config *fxconfig.Config
	logger *fxlogger.Logger
}

func NewTestService(config *fxconfig.Config, logger *fxlogger.Logger) *TestService {
	return &TestService{
		config: config,
		logger: logger,
	}
}

func (s *TestService) Test(c echo.Context) string {

	fxlogger.Ctx(c.Request().Context()).Info().Msg("[lecho-ctx] called TestService::Test()")

	return strings.ToUpper(s.config.AppName)
}
