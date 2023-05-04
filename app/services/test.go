package services

import (
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"strings"
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

func (s *TestService) Test() string {
	s.logger.Info().Msg("called TestService::Test()")

	return strings.ToUpper(s.config.AppName)
}
