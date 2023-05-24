package post

import (
	"errors"
	"fmt"
	"github.com/ekkinox/fx-template/app/repository"
	"gorm.io/gorm"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type DeletePostHandler struct {
	repository *repository.PostRepository
}

func NewDeletePostHandler(repository *repository.PostRepository) *DeletePostHandler {
	return &DeletePostHandler{
		repository: repository,
	}
}

func (h *DeletePostHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			msg := fmt.Sprintf("invalid id: %v", err)
			c.Logger().Error(msg)

			return echo.NewHTTPError(http.StatusBadRequest, msg)
		}

		post, err := h.repository.Find(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				msg := fmt.Sprintf("cannot get post with id %d: %v", id, err)
				c.Logger().Error(msg)

				return echo.NewHTTPError(http.StatusNotFound, msg)
			}

			c.Logger().Errorf("cannot get post: %v", err)
			return err
		}

		err = h.repository.Delete(post)
		if err != nil {
			c.Logger().Errorf("cannot delete post: %v", err)
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}
