package fxauthentication

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/labstack/echo/v4"
)

const AuthenticationContextKey = "_authentication_context"

func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			authHeader := c.Request().Header.Get(echo.HeaderAuthorization)

			authHeaderParts := strings.SplitN(authHeader, " ", 2)
			if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
				c.Logger().Errorf("invalid authorization header %s", authHeader)

				return next(c)
			}

			authTokenParts := strings.Split(authHeaderParts[1], ".")

			var encodedAuthContext string
			if len(authTokenParts) == 3 {
				encodedAuthContext = authTokenParts[1]
			} else if len(authTokenParts) == 1 {
				encodedAuthContext = authTokenParts[0]
			} else {
				c.Logger().Errorf("authorization header does not contain valid jwt: %s", authHeaderParts[1])

				return next(c)
			}

			base64DecodedAuthContext, err := base64.RawStdEncoding.DecodeString(encodedAuthContext)
			if err != nil {
				c.Logger().Errorf("cannot base64 decode authorization jwt payload %s: %v", encodedAuthContext, err)

				return next(c)
			}

			var authContext AuthenticationContext
			err = json.Unmarshal(base64DecodedAuthContext, &authContext)
			if err != nil {
				c.Logger().Errorf("cannot json unmarshall authorization jwt payload %s: %v", base64DecodedAuthContext, err)

				return next(c)
			}

			c.Set(AuthenticationContextKey, &authContext)

			return next(c)
		}
	}
}
