# Go Analytics Platform

A multi-service real-time analytics platform built with Go, featuring event ingestion, asynchronous processing, and a web dashboard.

## 🏗️ Architecture

This project implements a microservices architecture with:

- **API Service**: REST API for event ingestion and stats retrieval + Web UI
- **Processor Service**: Kafka consumer that processes events and stores them in PostgreSQL
- **Apache Kafka**: Message queue for asynchronous event processing
- **PostgreSQL**: Persistent storage for analytics events
- **Web Dashboard**: Real-time analytics visualization

## 🚀 Technologies

- **Go 1.23.4**
- **Gin** - HTTP web framework
- **Kafka (Confluent)** - Event streaming
- **PostgreSQL 15** - Database
- **sqlc** - Type-safe SQL query generator
- **Docker & Docker Compose** - Containerization

## 📁 Project Structure

```
go-analytics/
├── api/                      # API service
│   ├── handler.go           # HTTP handlers
│   ├── main.go              # API server entry point
│   ├── Dockerfile
│   └── ui/static/           # Web UI files
├── processor/               # Event processor service
│   ├── main.go              # Processor entry point
│   └── Dockerfile
├── db/
│   ├── migration/           # Database migrations
│   ├── sql/                 # SQL queries
│   └── sqlc/                # Generated Go code from sqlc
├── pkg/
│   ├── config/              # Configuration management
│   └── types/               # Shared types
├── docker-compose.yml       # Multi-container orchestration
└── go.mod                   # Go module dependencies
```

## 🛠️ Setup & Installation

### Prerequisites

- Docker & Docker Compose
- Go 1.23+ (for local development)

### Quick Start

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd go-analytics
   ```

2. **Start all services**
   ```bash
   docker-compose up --build
   ```

3. **Access the application**
   - Web UI: http://localhost:8080
   - API: http://localhost:8080/events (POST)
   - Stats: http://localhost:8080/stats (GET)

## 📊 API Endpoints

### POST `/events` - Ingest Event
Send analytics events to the platform.

**Request:**
```json
{
  "site_id": "site-abc-123",
  "event_type": "page_view",
  "path": "/home",
  "user_id": "user-001",
  "timestamp": "2025-11-15T05:00:00Z"
}
```

**Response:**
```json
{
  "message": "Event accepted"
}
```

### GET `/stats` - Retrieve Statistics
Query aggregated analytics for a site and date.

**Parameters:**
- `site_id` (required): Site identifier
- `date` (required): Date in YYYY-MM-DD format

**Example:**
```
GET /stats?site_id=site-abc-123&date=2025-11-15
```

**Response:**
```json
{
  "site_id": "site-abc-123",
  "date": "2025-11-15",
  "total_views": 150,
  "unique_users": 42,
  "top_paths": [
    {
      "path": "/home",
      "views": 50
    },
    {
      "path": "/products",
      "views": 35
    }
  ]
}
```

## 🧪 Testing

### Using Postman

1. **Import requests** or create new ones:
   - POST http://localhost:8080/events
   - GET http://localhost:8080/stats?site_id=site-abc-123&date=2025-11-15

2. **Send test events** with JSON body as shown above

### Using PowerShell

```powershell
# Send an event
$body = @{
    site_id = "site-abc-123"
    event_type = "page_view"
    path = "/home"
    user_id = "user-001"
    timestamp = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ss.fffZ")
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/events" -Method POST -Body $body -ContentType "application/json"

# Get stats
Invoke-RestMethod -Uri "http://localhost:8080/stats?site_id=site-abc-123&date=2025-11-15"
```

## 🔧 Local Development

### Run Infrastructure Only

```bash
docker-compose up -d zookeeper kafka postgres migrator
```

### Set Environment Variables

```bash
export DATABASE_URL="postgres://admin:secret@localhost:5432/analytics?sslmode=disable"
export KAFKA_BROKER_URLS="localhost:9092"
export KAFKA_TOPIC="events"
export KAFKA_GROUP_ID="processor-group"
```

### Run Services Locally

```bash
# Terminal 1 - API Service
go run ./api

# Terminal 2 - Processor Service
go run ./processor
```

## 🗄️ Database Schema

```sql
CREATE TABLE "events" (
  "id" SERIAL PRIMARY KEY,
  "site_id" varchar(100) NOT NULL,
  "event_type" varchar(50) NOT NULL,
  "path" text,
  "user_id" varchar(100),
  "timestamp" timestamptz NOT NULL
);
```

**Indexes:**
- `idx_events_site_id_timestamp`
- `idx_events_user_id`
- `idx_events_site_id_event_type_timestamp`

## 📈 Features

- ✅ Real-time event ingestion via REST API
- ✅ Asynchronous event processing with Kafka
- ✅ Type-safe database queries with sqlc
- ✅ Aggregated analytics (total views, unique users, top paths)
- ✅ Web-based dashboard for data visualization
- ✅ Dockerized microservices architecture
- ✅ Database migrations support
- ✅ Production-ready error handling

## 🐛 Troubleshooting

### Check Service Logs

```bash
docker-compose logs -f api
docker-compose logs -f processor
```

### Verify Database

```bash
docker exec postgres psql -U admin -d analytics -c "SELECT COUNT(*) FROM events;"
```

### Reset Everything

```bash
docker-compose down -v
docker-compose up --build
```

## 📝 License

MIT License - feel free to use this project for learning and development.

## 👤 Author

Created as a learning project for building microservices with Go.

## 🤝 Contributing

Contributions, issues, and feature requests are welcome!
