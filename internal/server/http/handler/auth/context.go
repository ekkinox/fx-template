package auth

import (
	"net/http"

	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
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

		authCtx, err := fxhttpserver.CtxAuthenticationContext(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid auth context")
		}

		return c.JSON(http.StatusOK, echo.Map{
			"full_context": authCtx,
			"analysis": echo.Map{
				"is_user":   authCtx.IsUserEntity(),
				"is_admin":  authCtx.IsAdminEntity(),
				"uuid":      authCtx.Uuid,
				"client_id": authCtx.ClientId,
				"idp":       authCtx.IdentityProviderType.String(),
				"entity":    authCtx.EntityType().String(),
			},
		})
	}
}
