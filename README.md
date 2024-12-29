# go-YTFetch

A Go-based REST API service that fetches and stores YouTube videos based on search queries. The service implements cursor-based pagination and efficient API key rotation, featuring a modern React frontend with responsive design.

## Features

- YouTube video search and storage
- Cursor-based pagination for efficient data retrieval
- Round-robin API key rotation
- PostgreSQL database integration
- RESTful API endpoints
- Background job for periodic video fetching
- Responsive React frontend with Tailwind CSS
- Error handling and loading states
- Configurable cron schedules for data fetching
- Optimized database queries and indexing
- TypeScript for enhanced type safety
- API response caching support
- CORS and security configurations

## Tech Stack

### Backend
- Go 1.21+
- Gin Web Framework
- PostgreSQL
- YouTube Data API v3

### Frontend
- React 18+
- TypeScript
- Tailwind CSS
- Vite


## Installation

1. Clone the repository:
```bash
git clone https://github.com/palSagnik/go-YTFetch.git
cd go-YTFetch
```

2. Install dependencies:
```bash
go mod download

cd frontend
npm install
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

2. The Frontend will be available at `http://localhost:5173`

## API Endpoints

### Fetch videos using YT API
```
GET /api/fetch
```
Query Parameters:
- `query`: Search query to fetch results
- `max_results` (optional): Number of results per query
- `published_after`: Date after which it was published


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

## Screenshots
### First results after fetching from database
<img width="1280" alt="Screenshot 2024-12-29 at 2 46 12 PM" src="https://github.com/user-attachments/assets/0ab80159-9a91-439b-baa0-3f3ca8a62f9b" />

### Pagination implementation and next results
<img width="1275" alt="Screenshot 2024-12-29 at 2 45 00 PM" src="https://github.com/user-attachments/assets/cee52e42-b22b-4951-8ba9-b36a8c3fab92" />

### Buttons to navigate pages
<img width="1280" alt="Screenshot 2024-12-29 at 3 52 15 PM" src="https://github.com/user-attachments/assets/3769eaf5-1a5d-4c27-a6be-340028730365" />


## License

This project is licensed under the MIT License - see the LICENSE file for details.
