package services

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"strings"

	"github.com/ekkinox/fx-template/modules/fxconfig"
)

type TestService struct {
	config *fxconfig.Config
}

func NewTestService(config *fxconfig.Config) *TestService {
	return &TestService{
		config: config,
	}
}

func (s *TestService) Test(c echo.Context) (string, error) {

	c.Logger().Info("service TestService invoked")

	if s.config.GetBool("APP_SHOULD_FAIL") {
		e := errors.New("failure")

		c.Logger().Errorf("app was configured to fail: %v", e)
		return "", e
	}

	return fmt.Sprintf("app name is: %s", strings.ToUpper(s.config.AppName())), nil
}
