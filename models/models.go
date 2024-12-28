package models

// video Item model containing fields required
type VideoItem struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	PublishedAt  string `json:"published_at"`
	ChannelTitle string `json:"channeltitle"`
	ThumbnailURL string `json:"thumbnailurl"`
}

// api response model
type VideoResponse struct {
	Status  string      `json:"status"`
	Data    []VideoItem `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// cursor pagination model
type PaginatedVideos struct {
	Videos     []VideoItem `json:"videos"`
	HasNext    bool        `json:"hasNext"`
	NextCursor string      `json:"nextCursor,omitempty"`
	TotalCount int64       `json:"totalCount"`
}

type PaginationQuery struct {
	Cursor    string
	Limit     int64
	SortField string
	SortOrder string
}
