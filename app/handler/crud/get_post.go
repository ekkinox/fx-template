package crud

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ekkinox/fx-template/app/repository"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type GetPostHandler struct {
	repository *repository.PostRepository
}

func NewGetPostHandler(repository *repository.PostRepository) *GetPostHandler {
	return &GetPostHandler{
		repository: repository,
	}
}

func (h *GetPostHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {

		c.Logger().Info("in get post handler")

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			msg := fmt.Sprintf("invalid id: %v", err)
			c.Logger().Error(msg)

			return echo.NewHTTPError(http.StatusBadRequest, msg)
		}

		post, err := h.repository.Find(c.Request().Context(), id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				msg := fmt.Sprintf("cannot get post with id %d: %v", id, err)
				c.Logger().Error(msg)

				return echo.NewHTTPError(http.StatusNotFound, msg)
			}

			c.Logger().Errorf("cannot get post: %v", err)
			return err
		}

		return c.JSON(http.StatusOK, post)
	}
}
