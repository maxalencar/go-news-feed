# News Feed API

A simple RESTful microservice that reads RSS feeds from public sources and stores them in MongoDB for future queries.

## Features

- **RSS Feed Aggregation**: Fetches articles from multiple news sources
- **MongoDB Storage**: Persists articles for efficient querying
- **Flexible Filtering**: Find articles by category, provider, date range, etc.
- **Pagination Support**: Efficient pagination with limit and page parameters
- **Sorting**: Sort results by any article property

## API Endpoints

### GET /load

Loads articles from RSS feeds into the database.

**Query Parameters:**
- `feedUrl` (optional): Specific feed URL to load. If not provided, loads from all default sources.

**Default Sources:**
- BBC News UK: https://feeds.bbci.co.uk/news/uk/rss.xml
- BBC News Technology: https://feeds.bbci.co.uk/news/technology/rss.xml
- Sky News UK: https://feeds.skynews.com/feeds/rss/uk.xml
- Sky News Technology: https://feeds.skynews.com/feeds/rss/technology.xml

**Example Requests:**
```bash
# Load from all default sources
curl -X GET http://localhost:8080/load

# Load from a specific feed
curl -X GET http://localhost:8080/load?feedUrl=https://feeds.skynews.com/feeds/rss/technology.xml
```

**Response Format:**
```json
[
  {
    "id": "https://www.bbc.co.uk/news/business-61644033",
    "title": "Could flat tyres soon be a thing of the past?",
    "description": "Airless tyres that do not puncture are getting close to market but some remain sceptical about them.",
    "link": "https://www.bbc.co.uk/news/business-61644033?at_medium=RSS&at_campaign=KARANGA",
    "source": {
      "category": "uk",
      "feedUrl": "https://feeds.bbci.co.uk/news/uk/rss.xml",
      "provider": "bbc"
    },
    "publishedDateTime": "2022-06-13T23:18:58Z"
  }
]
```

### GET /find

Searches and retrieves articles from the database based on filters.

**Query Parameters:**
| Parameter | Type   | Description |
|-----------|--------|-------------|
| category  | string | Article's category (uk, technology) |
| provider  | string | Article's provider (bbc, sky) |
| limit     | int    | Max results per page (default: 100, max: 1000) |
| page      | int    | Page index (default: 0) |
| sort      | string | Sort column (e.g., publishedDateTime) |
| order     | string | Sort order (asc or desc, default: asc) |

**Example Requests:**
```bash
# Find all articles
curl http://localhost:8080/find

# Find UK BBC articles sorted by date descending
curl http://localhost:8080/find?sort=publishedDateTime&order=desc&category=uk&provider=bbc
```

**Response Format:**
```json
{
  "criteria": {
    "category": "uk",
    "provider": "bbc",
    "sort": "publishedDateTime",
    "order": "desc"
  },
  "articles": [
    {
      "id": "https://www.bbc.co.uk/news/uk-62874346",
      "title": "King Charles III promises to follow Queen's selfless duty",
      "description": "The King hears condolences at Westminster before travelling to Edinburgh to mount a vigil for the Queen.",
      "link": "https://www.bbc.co.uk/news/uk-62874346?at_medium=RSS&at_campaign=KARANGA",
      "source": {
        "category": "uk",
        "feedUrl": "https://feeds.bbci.co.uk/news/uk/rss.xml",
        "provider": "bbc"
      },
      "publishedDateTime": "2022-09-12T12:47:57Z"
    }
  ],
  "total": 79
}
```

## Getting Started

### Prerequisites

You will need:
- Go 1.22 or higher
- MongoDB instance
- Docker (optional, for containerized deployment)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/maxalencar/go-news-feed.git
cd go-news-feed
```

2. Install dependencies:
```bash
go mod tidy
```

### Running the Application

#### Using Go directly:
```bash
# Set required environment variables
export MONGO_URI="mongodb://localhost:27017"
export MONGO_DATABASE="news"
export MONGO_COLLECTION="articles"
export PORT=8080

# Run the application
go run cmd/news/main.go
```

#### Using Docker Compose (recommended):
```bash
docker compose up --build
```

The API will be available at `http://localhost:8080`.

### Running Tests

```bash
go test ./... -v
```

## Project Structure

```
.
├── cmd/                 # Main application entry point
├── internal/            # Internal application logic
│   ├── news/           # News service components
│   │   ├── server.go   # Server initialization and startup
│   │   ├── endpoint.go # HTTP endpoints
│   │   ├── service.go  # Business logic
│   │   ├── repository.go # Data access layer
│   │   └── config.go   # Configuration handling
├── pkg/                 # Shared packages
│   └── model/          # Data models
├── script/              # Utility scripts
└── docker-compose.yaml  # Docker configuration
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `MONGO_URI` | MongoDB connection string | `mongodb://localhost:27017` |
| `MONGO_DATABASE` | Database name | `news` |
| `MONGO_COLLECTION` | Collection name | `articles` |
| `PORT` | Server port | `8080` |

## Performance Optimizations

This application includes several performance improvements:
- **Bulk Operations**: Efficient bulk write operations to reduce database round trips
- **Database Indexing**: Proper indexing for query performance  
- **N+1 Query Elimination**: Optimized article saving to avoid redundant database calls
- **Connection Pooling**: Efficient database connection management

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/your-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin feature/your-feature`)
5. Create a Pull Request

## License

This project is licensed under the MIT License.