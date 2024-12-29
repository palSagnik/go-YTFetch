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
GET /api/getVideos
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
├── backend
│   ├── Dockerfile
│   ├── config
│   │   ├── db.go
│   │   └── general.go
│   ├── database
│   │   ├── database.go
│   │   ├── queries.go
│   │   └── schemas.go
│   ├── handler
│   │   └── api.go
│   ├── models
│   │   └── models.go
│   ├── routes
│   │   └── routes.go
│   └── utils
│       └── utils.go
├── compose.yml
├── cron
│   ├── Dockerfile
│   └── cron.go
├── frontend
│   ├── eslint.config.js
│   ├── index.html
│   ├── package-lock.json
│   ├── package.json
│   ├── postcss.config.js
│   ├── public
│   │   └── vite.svg
│   ├── src
│   │   ├── App.css
│   │   ├── App.tsx
│   │   ├── assets
│   │   │   └── react.svg
│   │   ├── components
│   │   │   ├── VideoList.tsx
│   │   │   └── ui
│   │   │       ├── button.tsx
│   │   │       └── card.tsx
│   │   ├── index.css
│   │   ├── main.tsx
│   │   ├── types.ts
│   │   └── vite-env.d.ts
│   ├── tailwind.config.js
│   ├── tsconfig.app.json
│   ├── tsconfig.json
│   ├── tsconfig.node.json
│   └── vite.config.ts
├── go.mod
├── go.sum
├── main.go
└── scripts
    ├── init.sh
    ├── reset.sh
    └── run.sh
```

## Design Choices

### Cursor-Based Pagination
- Chosen over offset pagination due to better performance with large datasets
- Prevents the "skipping rows" problem when data is added/removed between pages
- More efficient for our use case where we're mostly showing recent videos
- Handles concurrent updates better than offset pagination
- Better memory utilization on database queries
```go
videos, err := database.QueryVideosWithCursor(c, models.PaginationQuery{
    Limit:      10,
    NextCursor: "cursor_value",
})
```

### API Key Rotation
- Round-robin implementation based on quota chosen for reliability
- Thread-safe implementation using mutex
```go
keyManager, err := utils.NewAPIKeyManager([]string{
    os.Getenv("YOUTUBE_API_KEY_1"),
    os.Getenv("YOUTUBE_API_KEY_2"),
    os.Getenv("YOUTUBE_API_KEY_3"),
})
```
### Frontend Architecture
- React with TypeScript for type safety
- Tailwind CSS for utility-first styling
- Grid layout for better mobile experience

### Database Schema
- Optimized indexes for common query patterns
- Timestamp fields for better data tracking
- Constraint choices to maintain data integrity


## License

This project is licensed under the MIT License - see the LICENSE file for details.