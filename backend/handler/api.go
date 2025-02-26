package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/palSagnik/go-YTFetch.git/backend/config"
	"github.com/palSagnik/go-YTFetch.git/backend/database"
	"github.com/palSagnik/go-YTFetch.git/backend/models"
	"github.com/palSagnik/go-YTFetch.git/backend/utils"
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
	if maxStr := c.Query("max_results"); maxStr != "" {
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
		publishedAfter = config.DEFAULT_PUBLISHED_AFTER
	}

	// apikey parameter for rotation
	// if apikey not provided default to APIKEY_3
	// you can use any apikey which is not being used in rotation for apikeymanager
	// using 3 here due to shortage of apikeys
	apiKey := c.Query("api_key")
	if apiKey == "" {
		apiKey = config.YOUTUBE_APIKEY_3
	}

	// calling the search function with params
	videos, err := utils.SearchYoutubevideos(apiKey, query, publishedAfter, maxResults)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.VideoResponse{
			Status: "error",
			Message: "failed to search videos: " + err.Error(),
		})
		return
	}

	if err := database.InsertVideoDetails(c, videos); err != nil {
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

// querying database for videos
func GetVideos(c *gin.Context) {

	// taking query parameters
	// setting limit default value
	limit := int64(config.DEFAULT_PAGE_LIMIT)
	if limitStr := c.Query("limit"); limitStr != "" {
		var err error
		limit, err = strconv.ParseInt(limitStr, 10, 64)
		if err != nil || limit <= 0 {
			c.JSON(http.StatusBadRequest, models.VideoResponse{
                Status:  "error",
                Message: "invalid limit parameter",
            })
            return
		}
	}

	sortField := c.Query("field")
	SortOrder := c.Query("order")
	nextCursor := c.Query("next_cursor")
	
	paginationQuery := models.PaginationQuery{
		Limit: limit,
		SortField: sortField,
		SortOrder: SortOrder,
		NextCursor: nextCursor,
	}

	paginatedVideos, err := database.QueryVideosWithCursor(c, paginationQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.VideoResponse{
			Status: "error",
			Message: "error in querying database: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.VideoResponse{
		Status: "success",
		Data: paginatedVideos.Videos,
		Pagination: models.PaginationInfo{
			HasNext: paginatedVideos.HasNext,
			NextCursor: paginatedVideos.NextCursor,
			TotalCount: paginatedVideos.TotalCount,
		},
	})
}