package crud

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ekkinox/fx-template/app/model"
	"github.com/ekkinox/fx-template/app/repository"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UpdatePostHandler struct {
	repository *repository.PostRepository
}

func NewUpdatePostHandler(repository *repository.PostRepository) *UpdatePostHandler {
	return &UpdatePostHandler{
		repository: repository,
	}
}

func (h *UpdatePostHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {

		c.Logger().Info("in update post handler")

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

		update := new(model.Post)
		if err = c.Bind(update); err != nil {
			c.Logger().Errorf("cannot bind post updates: %v", err)
			return err
		}

		err = h.repository.Update(c.Request().Context(), post, update)
		if err != nil {
			c.Logger().Errorf("cannot update post: %v", err)
			return err
		}

		return c.JSON(http.StatusOK, post)
	}
}
