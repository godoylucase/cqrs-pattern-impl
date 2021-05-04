package api

import (
	"github.com/gin-gonic/gin"
)

func Configure(h *Handler, router *gin.Engine) *gin.Engine {
	v1 := router.Group("/v1")
	{
		v1.GET("/articles-by-global-hash-tags", h.GetArticleByGlobalHashTags)
		//v1.POST("/articles/:articleID/paragraph")
		//v1.PUT("/articles/:articleID")
		//v1.PUT("/articles/:articleID/paragraph/:paraghaphID")
		//v1.DELETE("/articles/:articleID")
		//v1.DELETE("/articles/:articleID/paragraph/:paraghaphID")
	}

	return router
}
