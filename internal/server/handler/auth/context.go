package auth

import (
	"net/http"

	"github.com/ekkinox/fx-template/modules/fxauthentication"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/labstack/echo/v4"
)

type AuthContextHandler struct {
	config *fxconfig.Config
}

func NewAuthContextHandler(config *fxconfig.Config) *AuthContextHandler {
	return &AuthContextHandler{
		config: config,
	}
}

func (h *AuthContextHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {

		authContext := c.Get(fxauthentication.AuthenticationContextKey).(*fxauthentication.AuthenticationContext)

		return c.JSON(http.StatusOK, authContext)
	}
}
