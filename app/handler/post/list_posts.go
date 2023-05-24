package post

import (
	"github.com/ekkinox/fx-template/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ListPostsHandler struct {
	repository *repository.PostRepository
}

func NewListPostsHandler(repository *repository.PostRepository) *ListPostsHandler {
	return &ListPostsHandler{
		repository: repository,
	}
}

func (h *ListPostsHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {

		posts, err := h.repository.FindAll()
		if err != nil {
			c.Logger().Errorf("cannot list posts: %v", err)
			return err
		}

		return c.JSON(http.StatusOK, posts)
	}
}
