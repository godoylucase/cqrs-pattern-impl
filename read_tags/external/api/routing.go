package api

import (
	"github.com/gin-gonic/gin"
)

func Configure(h *Handler, router *gin.Engine) *gin.Engine {
	v1 := router.Group("/v1")
	{
		v1.GET("/articles-by-global-hash-tags", h.GetArticleByGlobalHashTags)
		v1.GET("/user-articles-by-source-url", h.GetUsersByArticle)
	}

	return router
}
