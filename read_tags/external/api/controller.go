package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/godoylucase/read_tags/business/dto"
)

type service interface {
	GetArticleByGlobalTags(globalHashTags []string) (dto.ArticleByGlobalHashTagRead, error)
	GetUsersByArticle(articleID string) ([]dto.UserByArticle, error)
}

type Handler struct {
	service service
}

func NewHandler(service service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetArticleByGlobalHashTags(c *gin.Context) {
	values := c.Request.URL.Query()
	ghts, ok := values["global_hash_tags"]
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	args := strings.Split(ghts[0], ",")
	aghts, err := h.service.GetArticleByGlobalTags(args)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, aghts)
}
