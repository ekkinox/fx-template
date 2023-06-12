package crud

import (
	"net/http"

	"github.com/ekkinox/fx-template/internal/app/repository"
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

		c.Logger().Info("in list posts handler")

		posts, err := h.repository.FindAll(c.Request().Context())
		if err != nil {
			c.Logger().Errorf("cannot list posts: %v", err)
			return err
		}

		return c.JSON(http.StatusOK, posts)
	}
}
