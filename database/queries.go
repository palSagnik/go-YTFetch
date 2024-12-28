package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/palSagnik/go-YTFetch.git/models"
)


func InsertVideoDetails (c *gin.Context, videoItems []models.VideoItem) error {

	if len(videoItems) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(videoItems))
	valueArgs := make([]interface{}, 0, len(videoItems) * 6)

	for i, video := range videoItems {

		// creating placeholders for each video
		placeholder := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)
		valueStrings = append(valueStrings, placeholder)

		publishedAt, err := time.Parse(time.RFC3339, video.PublishedAt)
		if err != nil {
            return fmt.Errorf("error parsing published_at for video %s: %v", video.ID, err)
        }
		
		valueArgs = append(valueArgs, 
		video.ID,
		video.Title,
		video.Description,
		publishedAt,
		video.ChannelTitle,
		video.ThumbnailURL,
		)
	}
	
	query := fmt.Sprintf(`INSERT INTO videos (
			video_id, title, description, published_at, channel_title, thumbnail_url)
			VALUES %s
			ON CONFLICT (video_id) DO UPDATE SET
			title = EXCLUDED.title,
            description = EXCLUDED.description,
            published_at = EXCLUDED.published_at,
            channel_title = EXCLUDED.channel_title,
            thumbnail_url = EXCLUDED.thumbnail_url`, 
			strings.Join(valueStrings, ","))
	
	tx, err := DB.BeginTx(c, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	_, err = tx.ExecContext(c, query, valueArgs...)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("error executing bulk insert: %v", err)
    }

	if err = tx.Commit(); err != nil {
        return fmt.Errorf("error committing transaction: %v", err)
    }

	return nil
}