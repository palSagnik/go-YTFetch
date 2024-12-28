package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/palSagnik/go-YTFetch.git/config"
	"github.com/palSagnik/go-YTFetch.git/database"
	"github.com/palSagnik/go-YTFetch.git/models"
	"github.com/palSagnik/go-YTFetch.git/utils"
)

func YTFetchApi(c *gin.Context) {

	// get query parameters
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, models.VideoResponse{
            Status:  "error",
            Message: "request parameter 'query' is required",
        })

		return
	}

	// default value for maxresult if no parameter is provided
	maxResults := int64(config.DEFAULT_VIDEO_FETCH_LIMIT) 
	if maxStr := c.Query("maxresults"); maxStr != "" {
		var err error
		maxResults, err = strconv.ParseInt(maxStr, 10, 64)
		if err != nil || maxResults <= 0 {
			c.JSON(http.StatusBadRequest, models.VideoResponse{
                Status:  "error",
                Message: "invalid 'max_results' parameter",
            })
            return
		}
	}
	
	publishedAfter := c.Query("published_after")
	if publishedAfter == "" {
		c.JSON(http.StatusBadRequest, models.VideoResponse{
			Status: "error",
			Message: "request parameter 'published_after' is required",
		})
	}


	// calling the search function with params
	videos, err := utils.SearchYoutubevideos(config.YOUTUBE_APIKEY, query, publishedAfter, maxResults)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.VideoResponse{
			Status: "error",
			Message: "failed to search videos: " + err.Error(),
		})
		return
	}

	var videoItems []models.VideoItem
	for _, video := range videos{
		videoItems = append(videoItems, models.VideoItem{
			ID: video.ID,
			Title: video.Title,
			Description: video.Description,
			PublishedAt: video.PublishedAt,
			ChannelTitle: video.ChannelTitle,
			ThumbnailURL: video.ThumbnailURL,
		})
	}

	if err := database.InsertVideoDetails(c, videoItems); err != nil {
		c.JSON(http.StatusInternalServerError, models.VideoResponse{
			Status: "error",
			Message: "failed to store videos in database: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.VideoResponse{
		Status: "success",
	})
}