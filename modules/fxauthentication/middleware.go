package fxauthentication

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const AuthenticationContextKey = "_authentication_context"

func Middleware() echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
			authHeaderParts := strings.Split(authHeader, ".")

			var encodedAuthContext string
			if len(authHeaderParts) == 3 {
				encodedAuthContext = authHeaderParts[1]
			} else if len(authHeaderParts) == 1 {
				encodedAuthContext = authHeaderParts[0]
			} else {
				c.Logger().Errorf("authorization header %s is invalid", authHeader)

				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			padding := len(encodedAuthContext) % 4
			if padding > 0 {
				encodedAuthContext += strings.Repeat("=", 4-padding)
			}

			base64DecodedAuthContext, err := base64.URLEncoding.DecodeString(encodedAuthContext)
			if err != nil {
				c.Logger().Errorf("cannot base64 decode authorization payload %s: %v", encodedAuthContext, err)

				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			var authContext AuthenticationContext
			err = json.Unmarshal(base64DecodedAuthContext, &authContext)
			if err != nil {
				c.Logger().Errorf("cannot json unmarshall authorization payload %s: %v", base64DecodedAuthContext, err)

				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			c.Set(AuthenticationContextKey, &authContext)

			return next(c)
		}
	}
}
