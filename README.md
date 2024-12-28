# go-YTFetch

A Go-based REST API service that fetches and stores YouTube videos based on search queries. The service implements cursor-based pagination and efficient API key rotation.

## Features

- YouTube video search and storage
- Cursor-based pagination for efficient data retrieval
- Round-robin API key rotation
- PostgreSQL database integration
- RESTful API endpoints
- Background job for periodic video fetching

## Tech Stack

- Go 1.21+
- Gin Web Framework
- PostgreSQL
- YouTube Data API v3

## Prerequisites

- Go 1.21 or higher
- PostgreSQL
- YouTube API key(s)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/palSagnik/go-YTFetch.git
cd go-YTFetch
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
```

Edit the `.env` file with your configuration:
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=your_database_name

# YouTube API Keys
YOUTUBE_API_KEY_1=your_first_api_key
YOUTUBE_API_KEY_2=your_second_api_key
YOUTUBE_API_KEY_3=your_third_api_key

# Server Configuration
PORT=8080
```


## Running the Application

1. Start the server:
```bash
scripts/run.sh
```

2. The API will be available at `http://localhost:8080`

## API Endpoints

### Get Videos
```
GET /api/videos
```

Query Parameters:
- `limit` (optional): Number of videos per page (default: 10)
- `next_cursor` (optional): Cursor for next page
- `field` (optional): Sort field (default: "published_at")
- `order` (optional): Sort order (default: "DESC")

Response:
```json
{
    "status": "success",
    "data": [
        {
            "id": "video_id",
            "title": "Video Title",
            "description": "Video Description",
            "published_at": "2024-12-29T10:00:00Z",
            "channel_title": "Channel Name",
            "thumbnail_url": "https://example.com/thumbnail.jpg"
        }
    ],
    "pagination": {
        "has_next": true,
        "next_cursor": "encoded_cursor_value",
        "count": 10,
        "total_count": 100
    }
}
```

## Project Structure

```
go-YTFetch/
├── README.md
├── config
│   ├── db.go
│   └── general.go
├── cron
│   └── cron.go
├── database
│   ├── database.go
│   ├── queries.go
│   └── schemas.go
├── go.mod
├── go.sum
├── handler
│   └── api.go
├── main.go
├── models
│   └── models.go
├── routes
│   └── routes.go
├── scripts
│   └── run.sh
├── test
└── utils
    └── utils.go
```

## Key Components

### API Key Manager
The service implements a round-robin algorithm for rotating YouTube API keys:
```go
keyManager, err := utils.NewAPIKeyManager([]string{
    os.Getenv("YOUTUBE_API_KEY_1"),
    os.Getenv("YOUTUBE_API_KEY_2"),
    os.Getenv("YOUTUBE_API_KEY_3"),
})
```

### Cursor-based Pagination
Efficient pagination implementation using cursors instead of offset:
```go
videos, err := database.QueryVideosWithCursor(c, models.PaginationQuery{
    Limit:      10,
    NextCursor: "cursor_value",
})
```

## Error Handling

The API returns structured error responses:
```json
{
    "status": "error",
    "message": "error description"
}
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.