package fxhttpserver

import (
	"net/http"

	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/go-errors/errors"
	"github.com/labstack/echo/v4"
)

func NewHttpServerErrorHandler(config *fxconfig.Config) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		he, ok := err.(*echo.HTTPError)
		if ok {
			if he.Internal != nil {
				if herr, ok := he.Internal.(*echo.HTTPError); ok {
					he = herr
				}
			}
		} else {
			he = &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
			}
		}

		code := he.Code
		message := he.Message
		if m, ok := he.Message.(string); ok {
			if config.AppDebug() {
				message = echo.Map{
					"message": m,
					"error":   err.Error(),
					"trace":   errors.New(err).ErrorStack(),
				}
			} else {
				message = echo.Map{"message": m}
			}
		}

		if c.Request().Method == http.MethodHead {
			err = c.NoContent(he.Code)
		} else {
			err = c.JSON(code, message)
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}
