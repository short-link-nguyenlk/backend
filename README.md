# Short Link API

A RESTful API service for creating and managing short URLs.

## Features

- Create short URLs from long URLs
- Retrieve original URLs from short codes
- Health check endpoint

## API Endpoints

### Create Short Link
```
POST /api/v1/shortlinks
Content-Type: application/json

{
    "original_url": "https://example.com/very/long/url"
}
```

### Get Short Link
```
GET /api/v1/shortlinks/{code}
```

### Health Check
```
GET /health
```

## Getting Started

1. Clone the repository
2. Run the server:
   ```bash
   go run cmd/api/main.go
   ```
3. The server will start on `http://localhost:8080`

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── handlers/
│   ├── models/
│   └── service/
└── README.md
```