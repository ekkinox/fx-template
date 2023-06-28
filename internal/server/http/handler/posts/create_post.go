package posts

import (
	"net/http"

	"github.com/ekkinox/fx-template/internal/model"
	"github.com/ekkinox/fx-template/internal/repository"
	"github.com/labstack/echo/v4"
)

type CreatePostHandler struct {
	repository *repository.PostRepository
}

func NewCreatePostHandler(repository *repository.PostRepository) *CreatePostHandler {
	return &CreatePostHandler{
		repository: repository,
	}
}

func (h *CreatePostHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {

		c.Logger().Info("in create post handler")

		post := new(model.Post)
		if err := c.Bind(post); err != nil {
			c.Logger().Errorf("cannot bind post: %v", err)
			return err
		}

		err := h.repository.Create(c.Request().Context(), post)
		if err != nil {
			c.Logger().Errorf("cannot create post: %v", err)
			return err
		}

		return c.JSON(http.StatusCreated, post)
	}
}
