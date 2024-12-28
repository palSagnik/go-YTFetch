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