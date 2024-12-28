package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/palSagnik/go-YTFetch.git/handler"
)

func SetUpRoutes(r *gin.Engine) {
	
	api := r.Group("/api")
	{
		api.GET("/fetch", handler.YTFetchApi)
		api.GET("/getVideos", handler.GetVideos)
	}

}