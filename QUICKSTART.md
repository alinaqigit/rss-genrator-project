# Quick Start Guide

Get the RSS Generator Project up and running in 5 minutes!

## 📋 Prerequisites

- Docker & Docker Compose
- Go 1.25.6+
- Git

## 🚀 Quick Start

### 1. Clone & Setup (2 minutes)

```bash
git clone https://github.com/alinaqigit/rss-generator-project.git
cd rss-generator-project

# Create .env file
cat > .env << EOF
PORT=8080
HOST=127.0.0.1
DATABASE_URL="postgres://postgres:postgres@localhost:5544/rss_generator?sslmode=disable"
EOF

# Install dependencies
go mod download
```

### 2. Start Database (1 minute)

```bash
docker-compose up -d postgres
# Wait ~5 seconds for database to be healthy
```

### 3. Run Migrations (1 minute)

```bash
# Install goose if needed
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
goose -dir ./sql/migrations postgres "$DATABASE_URL" up
```

### 4. Start Server (1 minute)

```bash
go run *.go
```

You should see:

```
Scraping on 10 goroutines every 1m0s duration
```

### ✅ Verify It Works

```bash
# In another terminal
curl http://localhost:8080/v1/healthz
# Should return: {}
```

**Success!** Server is running on `http://localhost:8080`

---

## 📡 Test Endpoints

### Create a User

```bash
curl -X POST http://localhost:8080/v1/user \
  -H "Content-Type: application/json" \
  -d '{"name": "bob"}'
```

Response:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-04-26T10:30:00Z",
  "updated_at": "2025-04-26T10:30:00Z",
  "name": "bob",
  "api_key": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
}
```

Save the `api_key`: You'll need it for authenticated requests.

### Get All Feeds

```bash
curl http://localhost:8080/v1/feed
```

### Create a Feed (Authenticated)

Replace `{api_key}` with your API key:

```bash
curl -X POST http://localhost:8080/v1/feed \
  -H "Authorization: ApiKey {api_key}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Hacker News",
    "url": "https://news.ycombinator.com/rss"
  }'
```

---

## 🛑 Stopping the Project

```bash
# Stop application
Ctrl + C

# Stop database
docker-compose down
```

---

## 🆘 Common Issues

### "DATABASE_URL is not set"

→ Create `.env` file (see step 1)

### "Failed to connect to database"

→ Check PostgreSQL is running: `docker-compose ps`

### Port already in use

→ Change PORT in `.env` or kill process: `lsof -ti:8080 | xargs kill`

### Permission denied on migrations

→ Run with `go` directly instead of downloaded binary

---

## 📚 Next Steps

- **Full Documentation**: Read [README.md](README.md)
- **API Details**: Check [API_REFERENCE.md](API_REFERENCE.md)
- **Architecture**: Study [ARCHITECTURE.md](ARCHITECTURE.md)
- **Development**: See [DEVELOPMENT.md](DEVELOPMENT.md)
- **Database**: Learn from [DATABASE.md](DATABASE.md)

---

## 💡 Pro Tips

### Save API Key to Environment Variable

```bash
export RSS_API_KEY="your_api_key_here"
curl -H "Authorization: ApiKey $RSS_API_KEY" http://localhost:8080/v1/user
```

### Monitor Scraper Activity

```bash
go run *.go 2>&1 | grep -i "scrap"
```

### Connect to Database Directly

```bash
psql $DATABASE_URL
# Inside psql:
SELECT * FROM users;
SELECT * FROM feeds;
SELECT COUNT(*) FROM posts;
```

### Rebuild Database

```bash
docker-compose down postgres
docker-compose up -d postgres
goose -dir ./sql/migrations postgres "$DATABASE_URL" up
```

---

**That's it!** You're ready to develop. Check out the full documentation for
more details.
