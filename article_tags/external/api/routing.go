package api

import (
	"github.com/gin-gonic/gin"
)

func Configure(h *Handler, router *gin.Engine) *gin.Engine {
	v1 := router.Group("/v1")
	{
		v1.POST("/articles", h.CreateArticle)
		v1.GET("/articles/:article_id", h.GetArticle)
		v1.PUT("/articles/:article_id", h.UpdateArticle)
		//v1.POST("/articles/:articleID/paragraph")
		//v1.PUT("/articles/:articleID/paragraph/:paraghaphID")
		//v1.DELETE("/articles/:articleID")
		//v1.DELETE("/articles/:articleID/paragraph/:paraghaphID")
	}

	return router
}
