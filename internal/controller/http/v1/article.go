package v1

import (
	"blog-backend/internal/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type articleRoutes struct {
	articleUseCase usecase.Article
}

func newArticleRoutes(g *echo.Group, articleUseCase usecase.Article) {
	r := &articleRoutes{
		articleUseCase: articleUseCase,
	}

	g.POST("/articles", r.create)
}

type createArticleInput struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Content     string `json:"content" validate:"required"`
}

func (r *articleRoutes) create(c echo.Context) error {
	var input createArticleInput

	err := BindAndValidate(c, &input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	articleID, err := r.articleUseCase.CreateArticle(c.Request().Context(), usecase.ArticleCreateArticleInput{
		AuthorID:    c.Get(userIDCtx).(uuid.UUID),
		Title:       input.Title,
		Description: input.Description,
		Content:     input.Content,
	})

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id": articleID,
	})
}
