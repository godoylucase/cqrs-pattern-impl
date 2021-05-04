package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godoylucase/articles_tags/business"
	"github.com/godoylucase/articles_tags/internal"
)

type articleService interface {
	Create(ctx context.Context, ba *business.BaseArticle) (string, error)
}

type Handler struct {
	as articleService
}

func NewHandler(as articleService) *Handler {
	return &Handler{as: as}
}

func (h *Handler) CreateArticle(c *gin.Context) {
	var article business.BaseArticle
	if err := c.BindJSON(&article); err != nil {
		c.Error(&internal.AppError{
			Cause: fmt.Errorf("error when parsing request body: %w", err),
			Type:  internal.ErrValueValidation,
		})
		return
	}

	// TODO validate body completion else appError -> 400

	id, err := h.as.Create(c, &article)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("x-resource_id", id)
	c.Writer.WriteHeader(http.StatusCreated)
}
