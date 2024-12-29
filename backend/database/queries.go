package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/palSagnik/go-YTFetch.git/backend/models"
	"github.com/palSagnik/go-YTFetch.git/backend/utils"
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

// video insertion in cron job
func InsertVideoDetailsCron (c context.Context, videoItems []models.VideoItem) error {

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


// using cursor pagination for querying videos
func QueryVideosWithCursor(c *gin.Context, query models.PaginationQuery) (*models.PaginatedVideos, error) {

	// default values
	if query.SortField == "" {
		query.SortField = "published_at"
	}
	if query.SortOrder == "" {
		query.SortOrder = "DESC"
	}

	// building the Base query
	baseQuery := `SELECT video_id, title, description, published_at, channel_title, thumbnail_url FROM videos`

	// cursor condition if provided
	var args []interface{}
	whereClause := ""
	if query.NextCursor != "" {
		decodedCursor, err := utils.DecodeCursor(query.NextCursor)
		if err != nil {
			return nil, fmt.Errorf("invalid cursor: %v", err)
		}

		// Adding where clause based on the sort order
		if strings.EqualFold(query.SortOrder, "DESC") {
			whereClause = fmt.Sprintf("WHERE %s < $1", query.SortField)
		} else {
			whereClause = fmt.Sprintf("WHERE %s > $1", query.SortField)
		}
		args = append(args, decodedCursor)
	}

	// full query
	// we will fetch one more than limit just to check if there are more results
	fullQuery := fmt.Sprintf("%s %s ORDER BY %s %s LIMIT %d", 
			baseQuery,
			whereClause,
			query.SortField,
			query.SortOrder,
			query.Limit + 1,
		)
	
	// executing the query
	rows, err := DB.QueryContext(c, fullQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("error in executing query: %v", err)
	}
	defer rows.Close()

	var videos []models.VideoItem
	for rows.Next() {
		var video models.VideoItem
		err := rows.Scan(
			&video.ID,
			&video.Title,
			&video.Description,
			&video.PublishedAt,
			&video.ChannelTitle,
			&video.ThumbnailURL,
		)
		if err != nil {
			return nil, fmt.Errorf("error in scanning rows: %v", err)
		}	
		videos = append(videos, video)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	// total count as an option if required
	var totalCount int64
	err = DB.QueryRowContext(c, "SELECT COUNT(*) FROM videos").Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("error getting total count: %v", err)
	}

	// to determine whether there are more results
	hasNext := false
	if int64(len(videos)) > query.Limit {
		hasNext = true
		videos = videos[:query.Limit] // remove the extra item
	}

	// generating the next cursor
	var nextCursor string
	if hasNext && len(videos) > 0 {
		lastVideo := videos[len(videos) - 1]
		
		var cursorValue interface{}
		switch query.SortField {
		case "published_at":
			cursorValue = lastVideo.PublishedAt
		case "id":
			cursorValue = lastVideo.ID
		default:
			cursorValue = lastVideo.PublishedAt
		}

		nextCursor = utils.EncodeCursor(cursorValue)
	}

	return &models.PaginatedVideos{
		Videos: videos,
		HasNext: hasNext,
		NextCursor: nextCursor,
		TotalCount: totalCount,
	}, nil
}