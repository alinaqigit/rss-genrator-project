# Project Structure Guide

Complete breakdown of the RSS Generator Project file organization and code
layout.

## 📁 Directory Tree

```
rss-generator-project/
│
├── 📄 README.md                    # Main project documentation
├── 📄 QUICKSTART.md                # 5-minute setup guide
├── 📄 API_REFERENCE.md             # Complete API documentation
├── 📄 ARCHITECTURE.md              # System design & architecture
├── 📄 DEVELOPMENT.md               # Development guide & best practices
├── 📄 DATABASE.md                  # Database schema & queries
├── 📄 PROJECT_STRUCTURE.md         # This file
├── 📄 go.mod                       # Go module definition
├── 📄 go.sum                       # Go dependency versions
├── 📄 docker-compose.yml           # Docker development setup
├── 📄 sqlc.yaml                    # sqlc code generation config
├── 📄 .env.example                 # Environment variables template
├── 📄 .gitignore                   # Git ignore rules
├── 🚀 main.go                      # Application entry point
│
├── 🌐 HTTP Handlers & Routing
│   ├── 📄 handler_user.go          # User management endpoints
│   ├── 📄 handler_feed.go          # Feed management endpoints
│   ├── 📄 handler_feedFollow.go    # Feed subscription endpoints
│   ├── 📄 handler-readiness.go     # Health check endpoint
│   └── 📄 handler_err.go           # Error handling demo
│
├── 🔐 Middleware & Utilities
│   ├── 📄 middleaare_auth.go       # API key authentication
│   ├── 📄 json.go                  # JSON response helpers
│   ├── 📄 models.go                # Data model conversions
│   ├── 📄 rss.go                   # RSS XML parsing
│   └── 📄 scraper.go               # Background RSS scraper
│
├── 📦 internal/
│   │
│   ├── 🔑 auth/
│   │   └── auth.go                 # API key extraction from headers
│   │
│   └── 💾 db/
│       ├── db.go                   # Database initialization
│       ├── models.go               # Database models (generated)
│       ├── users.sql.go            # User queries (generated)
│       ├── feeds.sql.go            # Feed queries (generated)
│       ├── feed_follows.sql.go     # Feed follow queries (generated)
│       └── posts.sql.go            # Post queries (generated)
│
├── 🗄️ sql/
│   │
│   ├── 📝 migrations/
│   │   ├── 001_users.sql           # Create users table
│   │   ├── 002_users_apikey.sql    # Add API key column
│   │   ├── 003_feeds.sql           # Create feeds table
│   │   ├── 004_feed_follows.sql    # Create feed_follows junction table
│   │   ├── 005_lastfetchedat.sql   # Add scraper timestamps
│   │   └── 006_posts.sql           # Create posts table
│   │
│   └── 🔍 queries/
│       ├── users.sql               # User SQL queries
│       ├── feeds.sql               # Feed SQL queries
│       ├── feed_follows.sql        # Feed follow SQL queries
│       └── posts.sql               # Post SQL queries
│
├── 📦 bin/
│   └── rss-generator               # Compiled binary (git-ignored)
│
└── 🔧 tmp/
    └── main                        # Development build (git-ignored)
```

---

## 📄 Root Level Files

### Configuration & Setup

| File                 | Purpose                             |
| -------------------- | ----------------------------------- |
| `go.mod`             | Module declaration & dependencies   |
| `go.sum`             | Dependency checksums (lock file)    |
| `sqlc.yaml`          | sqlc code generation configuration  |
| `docker-compose.yml` | Local development environment setup |
| `.env.example`       | Template for environment variables  |
| `.gitignore`         | Git ignore patterns                 |

### Documentation

| File                   | Purpose                             |
| ---------------------- | ----------------------------------- |
| `README.md`            | Main project overview (start here!) |
| `QUICKSTART.md`        | 5-minute setup guide                |
| `API_REFERENCE.md`     | Complete API documentation          |
| `ARCHITECTURE.md`      | System design & scalability         |
| `DEVELOPMENT.md`       | Development workflows & practices   |
| `DATABASE.md`          | Database schema & optimization      |
| `PROJECT_STRUCTURE.md` | This file                           |

---

## 🚀 Application Entry Point

### main.go

```go
main.go
├── package main
├── main() function
│   ├─ Load environment variables
│   ├─ Connect to PostgreSQL
│   ├─ Initialize API server
│   ├─ Configure CORS middleware
│   ├─ Register v1 routes:
│   │  ├─ /healthz
│   │  ├─ /user (POST, GET, DELETE)
│   │  ├─ /feed (POST, GET)
│   │  ├─ /feed-follow (POST, GET, DELETE)
│   │  └─ /posts (GET)
│   ├─ Start background scraper goroutine
│   └─ Start HTTP server
└── Total: ~100 lines
```

**Key Responsibilities:**

- Server initialization & configuration
- Database connection pooling
- Route registration
- Background task scheduling

---

## 🌐 HTTP Handlers (handler\_\*.go)

All handler files implement the same pattern:

```
handler_xyz.go
└── func (apiCfg *apiConfig) handlerNameMethod() {
    1. Parse request body
    2. Validate input
    3. Call database queries
    4. Transform data to API model
    5. Return JSON response
}
```

### handler_user.go

Manages user accounts:

- `handlerCreateUser`: POST /user - Register new account
- `handlerGetUserByAPI`: GET /user - Get user profile
- `handlerDeactivateUser`: DELETE /user - Delete account
- `handlerGetUserPosts`: GET /posts - Get user's posts (helper)

**Key Functions:**

```go
type User struct {
    ID        uuid.UUID
    Username  string
    ApiKey    string
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (apiCfg *apiConfig) handlerCreateUser(w, r)
func (apiCfg *apiConfig) handlerGetUserByAPI(w, r, user)
func (apiCfg *apiConfig) handlerDeactivateUser(w, r, user)
func (apiCfg *apiConfig) handlerGetUserPosts(w, r, user)
```

### handler_feed.go

Manages RSS feed sources:

- `handlerCreateUserFeed`: POST /feed - Add new feed
- `handlerGetAllFeeds`: GET /feed - List all feeds

**Key Functions:**

```go
type Feed struct {
    ID        uuid.UUID
    Name      string
    URL       string
    UserID    uuid.UUID
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (apiCfg *apiConfig) handlerCreateUserFeed(w, r, user)
func (apiCfg *apiConfig) handlerGetAllFeeds(w, r)
```

### handler_feedFollow.go

Manages feed subscriptions:

- `handlerCreateFeedFollow`: POST /feed-follow - Subscribe to feed
- `handlerGetFeedFollows`: GET /feed-follow - List user's subscriptions
- `handlerDeleteFeedFollows`: DELETE /feed-follow/{id} - Unsubscribe

**Key Functions:**

```go
type FeedFollow struct {
    ID        uuid.UUID
    UserID    uuid.UUID
    FeedID    uuid.UUID
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (apiCfg *apiConfig) handlerCreateFeedFollow(w, r, user)
func (apiCfg *apiConfig) handlerGetFeedFollows(w, r, user)
func (apiCfg *apiConfig) handlerDeleteFeedFollows(w, r, user)
```

### handler-readiness.go

Health check:

- `handlerReadiness`: GET /healthz - API health status

```go
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
    responseWithJson(w, 200, struct{}{})
}
```

### handler_err.go

Error demonstration (reference):

```go
func handlerErr(w http.ResponseWriter, r *http.Request) {
    responseWithError(w, 400, "Something went wrong")
}
```

---

## 🔐 Middleware & Utilities

### middleaare_auth.go

**Purpose:** API key authentication middleware

**Function:**

```go
type authHandler func(http.ResponseWriter, *http.Request, db.User)

func (apicfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
    // 1. Extract API key from Authorization header
    // 2. Query database for user
    // 3. Call handler with user context
    // 4. Return 403 if auth fails
}
```

**Usage:**

```go
// Protected endpoint
v1router.Post("/feed-follow",
    apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))

// Unprotected endpoint
v1router.Get("/feed", apiCfg.handlerGetAllFeeds)
```

### json.go

**Purpose:** JSON response utilities

**Functions:**

```go
func responseWithJson(w http.ResponseWriter, code int, payload any)
    └─ Marshal payload to JSON
    └─ Set Content-Type header
    └─ Write status code & body

func responseWithError(w http.ResponseWriter, code int, msg string)
    └─ Create error response object
    └─ Log 5XX errors
    └─ Call responseWithJson
```

**Example:**

```go
responseWithJson(w, 201, user)      // Success
responseWithError(w, 400, "Invalid") // Error
```

### models.go

**Purpose:** Data model definitions & conversions

**Structures:**

```go
type User struct                    // API User model
type Feed struct                    // API Feed model
type FeedFollow struct              // API FeedFollow model
type Post struct                    // API Post model

// Conversion functions
func database_user_to_User(db.User) User
func database_feed_to_Feed(db.Feed) Feed
func database_feedFollow_to_FeedFollow(db.FeedFollow) FeedFollow
func database_post_to_Post(db.Post) Post
func database_posts_to_Posts([]db.Post) []Post
```

**Pattern:** Separates database models (internal/db) from API response models

### rss.go

**Purpose:** RSS feed parsing

**Structures:**

```go
type RSSFeed struct {
    Channel struct {
        Title  string
        Link   string
        Description string
        Item   []RSSItem
    }
}

type RSSItem struct {
    Title       string
    Link        string
    Description string
    PubDate     string
}
```

**Function:**

```go
func urlToFeed(url string) (*RSSFeed, error)
    └─ HTTP GET with 10s timeout
    └─ Parse XML
    └─ Return structured feed
```

### scraper.go

**Purpose:** Background RSS feed scraping

**Functions:**

```go
func startScraping(db *db.Queries, concurrency int, interval time.Duration)
    └─ Initialize ticker
    └─ Fetch feeds to scrape
    └─ Launch concurrent goroutines

func scrapFeed(dbQueries *db.Queries, wg *sync.WaitGroup, feed db.Feed)
    └─ Fetch RSS from URL
    └─ Parse XML
    └─ Create post records
    └─ Update feed metadata
```

---

## 📦 internal/ Directory

### internal/auth/auth.go

**Purpose:** API key extraction from HTTP headers

**Function:**

```go
func GetAPIKey(headers http.Header) (string, error)
    └─ Extract Authorization header
    └─ Parse "ApiKey {key}" format
    └─ Return key or error
```

**Format:**

```
Authorization: ApiKey a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6
                     ^-- Exact match required
```

### internal/db/ Directory

Auto-generated by sqlc from SQL query files.

**Files:**

- `db.go` - Database connection initialization
- `models.go` - Generated database model structs
- `users.sql.go` - Generated user query functions
- `feeds.sql.go` - Generated feed query functions
- `feed_follows.sql.go` - Generated follow query functions
- `posts.sql.go` - Generated post query functions

**Example Generated Function:**

```go
// From sql/queries/users.sql
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
    row := q.queryRow(ctx, createUser, arg.ID, arg.Username)
    var i User
    err := row.Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt, &i.Username, &i.ApiKey)
    return i, err
}
```

---

## 🗄️ sql/ Directory

### sql/migrations/

Database schema evolution scripts using Goose format.

**Naming Convention:** NNN_description.sql (NNN = sequence number)

```
001_users.sql              # Initial user table
002_users_apikey.sql       # Add authentication
003_feeds.sql              # RSS feed tracking
004_feed_follows.sql       # Subscriptions
005_lastfetchedat.sql      # Scraper metadata
006_posts.sql              # Article storage
```

**Format:**

```sql
-- +goose up
-- Forward migration (apply new schema)
CREATE TABLE ...

-- +goose down
-- Rollback migration (revert schema)
DROP TABLE ...
```

### sql/queries/

SQL query definitions for sqlc code generation.

**Files:** `[resource].sql`

**Query Naming:** `-- name: QueryName :resultType`

```sql
-- name: CreateUser :one
-- Returns single row
INSERT INTO users (...) VALUES (...) RETURNING *;

-- name: GetAllFeeds :many
-- Returns multiple rows
SELECT * FROM feeds ORDER BY created_at;

-- name: DeletePost :exec
-- Executes, no return value
DELETE FROM posts WHERE id = $1;
```

---

## Development Patterns

### Adding New Feature

**Steps & Files:**

```
1. Identify feature scope
   └─ Affects database? → Add SQL migration & queries

2. Create SQL files
   └─ sql/migrations/NNN_feature.sql
   └─ sql/queries/feature.sql

3. Generate code
   └─ sqlc generate
   └─ Updates internal/db/feature.sql.go

4. Create handler
   └─ handler_feature.go
   └─ receiver methods on apiConfig

5. Add model conversions
   └─ models.go
   └─ database_X_to_X functions

6. Register routes
   └─ main.go
   └─ Add v1router.{Method}("/endpoint", handler)

7. Test endpoints
   └─ Manual curl tests
   └─ Add unit tests if needed
```

---

## Code Organization Principles

### 1. Separation of Concerns

- **SQL Logic**: `sql/queries/`
- **Database Layer**: `internal/db/`
- **HTTP Handlers**: `handler_*.go`
- **Authentication**: `middleweaare_*.go`
- **Data Transformation**: `models.go`
- **Business Logic**: `scraper.go`, `rss.go`

### 2. Error Handling Pattern

```go
// Always check errors
if err != nil {
    responseWithError(w, 400, fmt.Sprint("Description", err))
    return
}
```

### 3. Dependency Injection

```go
// apiConfig holds dependencies
type apiConfig struct {
    DB *db.Queries
}

// Methods use receiver for access
func (apiCfg *apiConfig) handler(...) {
    apiCfg.DB.QueryMethod(...)
}
```

### 4. Type Safety

- Database queries are type-safe (sqlc)
- JSON marshaling uses struct tags
- UUID generation via google/uuid package

---

## File Statistics

| Category   | Count  | Purpose                                  |
| ---------- | ------ | ---------------------------------------- |
| Handlers   | 5      | Route handlers (Users, Feeds, Follows)   |
| Middleware | 1      | Authentication middleware                |
| Utilities  | 4      | JSON, Models, RSS, Scraper               |
| Auth       | 1      | API key extraction                       |
| Database   | 4      | Models + generated code                  |
| Migrations | 6      | Schema evolution                         |
| Queries    | 4      | SQL query files                          |
| Config     | 4      | go.mod, docker-compose, .gitignore, sqlc |
| **Docs**   | **7**  | README, API, Architecture, etc.          |
| **Total**  | **36** | Documents + Code Files                   |

---

## Code Flow Example: Creating a Feed

```
HTTP Request
  POST /v1/feed

        ↓

Chi Router
  └─ Match route → Invoke middleware

        ↓

Middleware (middlewareAuth)
  └─ Extract API key from header
  └─ Query: users.sql.go::GetUserByAPIKey()
  └─ Get user object
  └─ Call handler with user

        ↓

Handler (handler_feed.go::handlerCreateUserFeed)
  └─ Parse JSON body
  └─ Create timestamp
  └─ Generate UUID
  └─ Call: feeds.sql.go::CreateFeed()
    └─ Execute: sql/queries/feeds.sql
    └─ INSERT INTO feeds
    └─ Return db.Feed

        ↓

Model Conversion (models.go)
  └─ Convert db.Feed → Feed (API model)
  └─ Extract needed fields

        ↓

Response (json.go::responseWithJson)
  └─ Marshal to JSON
  └─ Set headers
  └─ Write response (201 Created)

        ↓

HTTP Response
  200 Created with Feed JSON
```

---

## GitIgnore Patterns

```
/bin/                      # Compiled binaries
/tmp/                      # Temporary build files
.env                       # Environment variables (local)
*.log                      # Log files
.DS_Store                  # macOS files
.vscode/                   # Editor settings
dist/                      # Distribution directory
vendor/                    # Go vendor (if used)
*.prof                     # Profiling data
```

---

## Quick Reference Commands

```bash
# View structure
tree -L 2 -I 'vendor|bin'

# Count lines of code
find . -name "*.go" | xargs wc -l

# Build
go build -o bin/rss-generator

# Format
gofmt -s -w .

# Lint
golangci-lint run

# Test
go test ./... -v

# Generate DB code
sqlc generate

# Run migrations
goose -dir ./sql/migrations postgres "$DATABASE_URL" up
```

---

**Last Updated:** April 26, 2025
