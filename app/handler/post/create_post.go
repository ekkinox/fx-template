package post

import (
	"github.com/ekkinox/fx-template/app/model"
	"github.com/ekkinox/fx-template/app/repository"
	"net/http"

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

		post := new(model.Post)
		if err := c.Bind(post); err != nil {
			c.Logger().Errorf("cannot bind post: %v", err)
			return err
		}

		err := h.repository.Create(post)
		if err != nil {
			c.Logger().Errorf("cannot create post: %v", err)
			return err
		}

		return c.JSON(http.StatusCreated, post)
	}
}
