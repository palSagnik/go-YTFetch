package config

import (
	"os"
	"time"
)

var YOUTUBE_APIKEY_1 = os.Getenv("YOUTUBE_APIKEY_1")
var YOUTUBE_APIKEY_2 = os.Getenv("YOUTUBE_APIKEY_2")
var YOUTUBE_APIKEY_3 = os.Getenv("YOUTUBE_APIKEY_3")

var FETCH_YTAPI_AFTER = 15

var DEFAULT_PUBLISHED_AFTER = time.Now().UTC().AddDate(0, 0, -15).Truncate(24 * time.Hour).Format(time.RFC3339)
var DEFAULT_VIDEO_FETCH_LIMIT = 20
var DEFAULT_PAGE_LIMIT = 10
var APP_PORT = 9000