package fxauthenticationcontext

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const AuthenticationContextKey = "_authentication_context"

func Middleware(blockOnFailure bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
			authHeaderParts := strings.SplitN(authHeader, " ", 2)

			if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
				msg := "authorization header is not a token bearer"
				if blockOnFailure {
					c.Logger().Error(msg)

					return echo.NewHTTPError(http.StatusUnauthorized, msg)

				} else {
					c.Logger().Warnf("%s, skipping", msg)

					return next(c)
				}
			}

			authTokenParts := strings.Split(authHeaderParts[1], ".")

			var encodedAuthContext string
			if len(authTokenParts) == 3 {
				encodedAuthContext = authTokenParts[1]
			} else if len(authTokenParts) == 1 {
				encodedAuthContext = authTokenParts[0]
			} else {
				msg := fmt.Sprintf("authorization header does not contain valid jwt: %s", authHeaderParts[1])
				if blockOnFailure {
					c.Logger().Error(msg)

					return echo.NewHTTPError(http.StatusUnauthorized, msg)
				} else {
					c.Logger().Warnf("%s, skipping", msg)

					return next(c)
				}
			}

			base64DecodedAuthContext, err := base64.RawStdEncoding.DecodeString(encodedAuthContext)
			if err != nil {
				msg := fmt.Sprintf("cannot base64 decode authorization jwt payload %s: %v", encodedAuthContext, err)
				if blockOnFailure {
					c.Logger().Error(msg)

					return echo.NewHTTPError(http.StatusUnauthorized, msg)
				} else {
					c.Logger().Warnf("%s, skipping", msg)

					return next(c)
				}
			}

			var authContext AuthenticationContext
			err = json.Unmarshal(base64DecodedAuthContext, &authContext)
			if err != nil {
				msg := fmt.Sprintf("cannot json unmarshall authorization jwt payload %s: %v, skipping", base64DecodedAuthContext, err)
				if blockOnFailure {
					c.Logger().Error(msg)

					return echo.NewHTTPError(http.StatusUnauthorized, msg)
				} else {
					c.Logger().Warnf("%s, skipping", msg)

					return next(c)
				}
			}

			c.Set(AuthenticationContextKey, &authContext)

			return next(c)
		}
	}
}
