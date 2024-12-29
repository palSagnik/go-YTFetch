package utils

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/palSagnik/go-YTFetch.git/backend/models"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func SearchYoutubevideos(apiKey string, query string, publishedAfter string, maxResults int64) ([]models.VideoItem, error) {
	
	ctx := context.Background()
	youtube, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating YouTube client: %v", err)
	}

	// creating the search call
	call := youtube.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(maxResults).
		Order("date").
		Type("video").
		PublishedAfter(publishedAfter)

	// executing the search call
	resp, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("error executing search: %v", err)
	}

	var videos []models.VideoItem
	for _, item := range resp.Items {
		video := models.VideoItem{
			Title:        item.Snippet.Title,
			ID:           item.Id.VideoId,
			Description:  item.Snippet.Description,
			PublishedAt:  item.Snippet.PublishedAt,
			ChannelTitle: item.Snippet.ChannelTitle,
			ThumbnailURL: item.Snippet.Thumbnails.Default.Url,
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func EncodeCursor(value interface{}) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", value)))
}

func DecodeCursor(cursor string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}
