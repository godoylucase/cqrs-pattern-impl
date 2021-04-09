package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godoylucase/cqrs-pattern-impl/business"
	"github.com/godoylucase/cqrs-pattern-impl/internal"
)

type ArticleService interface {
	Create(ctx context.Context, ba *business.BaseArticle) (string, error)
}

type Handler struct {
	As ArticleService
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

	id, err := h.As.Create(c, &article)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("x-resource_id", id)
	c.Writer.WriteHeader(http.StatusCreated)
}
