# Database Documentation

Complete guide to the RSS Generator Project database schema, queries, and optimization strategies.

## Table of Contents

- [Database Overview](#database-overview)
- [Schema Design](#schema-design)
- [Tables](#tables)
- [Relationships](#relationships)
- [Indexes](#indexes)
- [Query Patterns](#query-patterns)
- [Optimization Tips](#optimization-tips)
- [Backup & Recovery](#backup--recovery)
- [Migration Guide](#migration-guide)

---

## Database Overview

### Technology Stack

- **DBMS**: PostgreSQL 16
- **Port**: 5544 (default in docker-compose)
- **Connection Method**: libpq driver
- **Query Management**: sqlc (type-safe compiled queries)
- **Migration Tool**: Goose

### Connection String Format

```
postgres://user:password@host:port/database?sslmode=disable
```

**Example:**
```
postgres://postgres:postgres@localhost:5544/rss_generator?sslmode=disable
```

### Connection Parameters

| Parameter | Value | Description |
|-----------|-------|-------------|
| User | `postgres` | Default PostgreSQL user (change in production) |
| Password | `postgres` | Default password (change in production) |
| Host | `localhost` or service name | Database server address |
| Port | `5544` (Docker) or `5432` (Standard) | PostgreSQL port |
| Database | `rss_generator` | Database name |
| SSL Mode | `disable` (dev) or `require` (prod) | SSL encryption requirement |

### Docker Exec

Access database from Docker:

```bash
# Connect to running container
docker-compose exec postgres psql -U postgres rss_generator

# Or with full connection string
docker-compose exec postgres psql "postgres://postgres:postgres@localhost:5432/rss_generator"
```

---

## Schema Design

### Design Principles

1. **Normalization**: Third normal form (3NF)
2. **Data Integrity**: Foreign keys with cascading deletes
3. **Audit Trail**: All tables have `created_at` and `updated_at`
4. **Uniqueness Constraints**: Enforced at database level
5. **Performance**: Strategic indexing on query columns

### Evolution Timeline

```
Migration 001: Create users table (base structure)
         ↓
Migration 002: Add API key column (authentication)
         ↓
Migration 003: Create feeds table (RSS sources)
         ↓
Migration 004: Create feed_follows table (user subscriptions)
         ↓
Migration 005: Add last_fetched_at to feeds (scraper tracking)
         ↓
Migration 006: Create posts table (RSS items)
```

---

## Tables

### 1. Users Table

Stores user account information.

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    username TEXT NOT NULL UNIQUE,
    api_key VARCHAR(255) UNIQUE NOT NULL DEFAULT (
        encode(sha256(random()::text::bytea), 'hex')
    )
);
```

**Columns:**

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Unique identifier (generated client-side) |
| `created_at` | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | Account creation timestamp |
| `updated_at` | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | Last modification timestamp |
| `username` | TEXT | NOT NULL, UNIQUE | Unique username (case-sensitive) |
| `api_key` | VARCHAR(255) | NOT NULL, UNIQUE | Auto-generated SHA256 hex API key |

**Indexes:**
- PRIMARY KEY: `id`
- UNIQUE: `username`, `api_key`

**Size Estimate**: ~500 bytes per row

**Example Data:**
```sql
SELECT * FROM users LIMIT 1;
```

```
                  id                  |       created_at       |       updated_at       | username |                api_key
--------------------------------------+------------------------+------------------------+----------+----------------------------------------
 550e8400-e29b-41d4-a716-446655440000 | 2025-04-26 10:30:00+00 | 2025-04-26 10:30:00+00 | alice    | a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6
```

---

### 2. Feeds Table

Stores RSS feed sources.

```sql
CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_fetched_at TIMESTAMPTZ
);
```

**Columns:**

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Unique feed identifier |
| `created_at` | TIMESTAMP | NOT NULL, DEFAULT NOW() | When feed was added |
| `updated_at` | TIMESTAMP | NOT NULL, DEFAULT NOW() | Last modification time |
| `name` | TEXT | NOT NULL | Human-readable feed name |
| `url` | TEXT | NOT NULL | RSS feed URL |
| `user_id` | UUID | NOT NULL, FK | Creator of this feed |
| `last_fetched_at` | TIMESTAMPTZ | DEFAULT NULL | Timestamp of last scrape |

**Indexes:**
- PRIMARY KEY: `id`
- FOREIGN KEY: `user_id` → `users(id)`
- Recommended: `last_fetched_at` (for scraper queries)

**Cascade Behavior**: When user deleted, all their feeds deleted

**Size Estimate**: ~1 KB per row

**Example Data:**
```sql
SELECT id, name, url, created_at, last_fetched_at FROM feeds LIMIT 1;
```

```
                  id                  |      name      |         url          |       created_at       |    last_fetched_at
--------------------------------------+----------------+----------------------+------------------------+------------------------
 660e8400-e29b-41d4-a716-446655440000 | Tech News      | https://ex.com/feed  | 2025-04-26 10:35:00    | 2025-04-26 10:40:25+00
```

---

### 3. Feed Follows Table

Junction table for user-feed subscriptions (many-to-many).

```sql
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    UNIQUE (user_id, feed_id)
);
```

**Columns:**

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Unique follow relationship ID |
| `created_at` | TIMESTAMPTZ | NOT NULL | When user started following |
| `updated_at` | TIMESTAMPTZ | NOT NULL | Last update |
| `user_id` | UUID | NOT NULL, FK | User ID |
| `feed_id` | UUID | NOT NULL, FK | Feed ID |
| UNIQUE | (user_id, feed_id) | - | Prevents duplicate follows |

**Indexes:**
- PRIMARY KEY: `id`
- FOREIGN KEY: `user_id` → `users(id)`
- FOREIGN KEY: `feed_id` → `feeds(id)`
- Unique Constraint: `(user_id, feed_id)`
- Recommended: `user_id`, `feed_id` (for lookups)

**Cascade Behavior**: 
- Delete user → delete all their follows
- Delete feed → delete all follow relationships

**Size Estimate**: ~150 bytes per row

**Example Data:**
```sql
SELECT uf.id, u.username, f.name 
  FROM feed_follows uf
  JOIN users u ON uf.user_id = u.id
  JOIN feeds f ON uf.feed_id = f.id
  LIMIT 1;
```

```
                  id                  | username |    name
--------------------------------------+----------+-----------
 880e8400-e29b-41d4-a716-446655440000 | alice    | Tech News
```

---

### 4. Posts Table

Stores individual RSS feed items (articles).

```sql
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    published_at TIMESTAMP NOT NULL,
    url TEXT NOT NULL UNIQUE,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE
);
```

**Columns:**

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Unique post ID |
| `created_at` | TIMESTAMP | NOT NULL | When added to database |
| `updated_at` | TIMESTAMP | NOT NULL | Last modification |
| `title` | TEXT | NOT NULL | Article title |
| `description` | TEXT | DEFAULT NULL | Article content/summary |
| `published_at` | TIMESTAMP | NOT NULL | Publication date from RSS |
| `url` | TEXT | NOT NULL, UNIQUE | Article URL (prevents duplicates) |
| `feed_id` | UUID | NOT NULL, FK | Source feed |

**Indexes:**
- PRIMARY KEY: `id`
- UNIQUE: `url`
- FOREIGN KEY: `feed_id` → `feeds(id)`
- Recommended: `feed_id`, `published_at DESC`

**Cascade Behavior**: Delete feed → delete all posts from that feed

**Size Estimate**: ~2-5 KB per row (varies by description length)

**Example Data:**
```sql
SELECT id, title, published_at, feed_id FROM posts ORDER BY published_at DESC LIMIT 1;
```

```
                  id                  |       title        |       published_at     |               feed_id
--------------------------------------+--------------------+------------------------+--------------------------------------
 aa0e8400-e29b-41d4-a716-446655440000 | Breaking News      | 2025-04-26 10:50:00+00 | 660e8400-e29b-41d4-a716-446655440000
```

---

## Relationships

### ER Diagram

```
┌─────────────────┐
│     users       │◄─────┐
├─────────────────┤      │
│ PK: id          │      │ (1:N)
│    username     │      │
│    api_key      │      │
└─────────────────┘      │
         │                │
         │ (1:N)          │
         │                │
┌─────────────────┐      │
│     feeds       │      │
├─────────────────┤      │
│ PK: id          │      │
│ FK: user_id ────┼──────┘
│    name         │
│    url          │
└─────────────────┘
         │
         │ (1:N)
         │
    ┌────┴────────────────────────────┐
    │                                 │
┌─────────────────┐         ┌─────────────────┐
│ feed_follows    │         │     posts       │
├─────────────────┤         ├─────────────────┤
│ PK: id          │         │ PK: id          │
│ FK: user_id ────┼─┐       │ FK: feed_id ────┼──┐
│ FK: feed_id ────┼─┼────────┤    title        │  │
│ UNIQUE(user,fd) │ │        │    url          │  │
└─────────────────┘ │        └─────────────────┘  │
                    │                              │
                    └──────────────────────────────┘
                         Many-to-many through
                        feed_follows & posts
```

### Relationship Types

**1. User → Feeds (1:N)**
- One user creates many feeds
- Foreign key: `feeds.user_id` → `users.id`
- Cascade: Delete user → delete all their feeds

**2. Feeds → Posts (1:N)**
- One feed has many posts
- Foreign key: `posts.feed_id` → `feeds.id`
- Cascade: Delete feed → delete all posts from feed

**3. Users ↔ Feeds (M:N via feed_follows)**
- One user can follow many feeds
- One feed can be followed by many users
- Junction table: `feed_follows`
- Cascade: Delete user/feed → delete follow relationship

### Query Patterns

**Get posts for user:**
```sql
SELECT p.* FROM posts p
  JOIN feed_follows uf ON p.feed_id = uf.feed_id
  WHERE uf.user_id = $1
  ORDER BY p.published_at DESC
  LIMIT 10;
```

**Get feeds for user:**
```sql
SELECT f.* FROM feeds f
  WHERE f.user_id = $1
  ORDER BY f.created_at DESC;
```

**Get user's followed feeds:**
```sql
SELECT f.* FROM feeds f
  JOIN feed_follows uf ON f.id = uf.feed_id
  WHERE uf.user_id = $1;
```

---

## Indexes

### Current Indexes

```sql
-- Primary keys (automatic)
CREATE INDEX idx_users_pk ON users(id);
CREATE INDEX idx_feeds_pk ON feeds(id);
CREATE INDEX idx_feed_follows_pk ON feed_follows(id);
CREATE INDEX idx_posts_pk ON posts(id);

-- Unique constraints (automatic)
CREATE UNIQUE INDEX idx_users_username ON users(username);
CREATE UNIQUE INDEX idx_users_api_key ON users(api_key);
CREATE UNIQUE INDEX idx_posts_url ON posts(url);
CREATE UNIQUE INDEX idx_feed_follows_unique ON feed_follows(user_id, feed_id);

-- Foreign keys
CREATE INDEX idx_feeds_user_id ON feeds(user_id);
CREATE INDEX idx_feed_follows_user_id ON feed_follows(user_id);
CREATE INDEX idx_feed_follows_feed_id ON feed_follows(feed_id);
CREATE INDEX idx_posts_feed_id ON posts(feed_id);
```

### Recommended Additional Indexes

Add these for production optimization:

```sql
-- Scraper query optimization
CREATE INDEX idx_feeds_last_fetched_at ON feeds(last_fetched_at);

-- Post queries
CREATE INDEX idx_posts_published_at_desc ON posts(published_at DESC);
CREATE INDEX idx_posts_feed_id_published_at ON posts(feed_id, published_at DESC);

-- Audit queries
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_posts_created_at ON posts(created_at);
```

### Creating Index

```bash
# Connect to database
psql $DATABASE_URL

# Create index
CREATE INDEX idx_feeds_last_fetched_at ON feeds(last_fetched_at);

# Verify
\d feeds
```

### Index Maintenance

```sql
-- Analyze index usage
SELECT schemaname, tablename, indexname, idx_scan
  FROM pg_stat_user_indexes
  ORDER BY idx_scan DESC;

-- Reindex if fragmented
REINDEX INDEX idx_posts_feed_id;

-- Monitor index size
SELECT indexrelname, pg_size_pretty(pg_relation_size(indexrelid)) as size
  FROM pg_stat_user_indexes
  ORDER BY pg_relation_size(indexrelid) DESC;
```

---

## Query Patterns

### User Queries

**Get user by API key:**
```sql
SELECT * FROM users WHERE api_key = $1;
```

**Create user:**
```sql
INSERT INTO users (id, username)
VALUES ($1, $2)
RETURNING *;
```

**Get user by ID:**
```sql
SELECT * FROM users WHERE id = $1;
```

**Delete user:**
```sql
DELETE FROM users WHERE api_key = $1;
```

### Feed Queries

**Create feed:**
```sql
INSERT INTO feeds (id, name, created_at, updated_at, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
```

**Get all feeds:**
```sql
SELECT * FROM feeds ORDER BY created_at DESC;
```

**Get feeds to scrape:**
```sql
SELECT * FROM feeds
  WHERE last_fetched_at IS NULL 
    OR last_fetched_at < NOW() - INTERVAL '1 minute'
  ORDER BY last_fetched_at ASC
  LIMIT $1;
```

**Update last fetched:**
```sql
UPDATE feeds
SET last_fetched_at = NOW()
WHERE id = $1
RETURNING *;
```

### Post Queries

**Create post:**
```sql
INSERT INTO posts 
  (id, created_at, updated_at, title, description, published_at, url, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (url) DO NOTHING
RETURNING *;
```

**Get posts for user:**
```sql
SELECT p.* FROM posts p
  JOIN feed_follows uf ON p.feed_id = uf.feed_id
  WHERE uf.user_id = $1
  ORDER BY p.published_at DESC
  LIMIT $2;
```

**Get recent posts:**
```sql
SELECT * FROM posts
  ORDER BY published_at DESC
  LIMIT $1;
```

### Feed Follow Queries

**Create follow:**
```sql
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
```

**Get user's follows:**
```sql
SELECT * FROM feed_follows
  WHERE user_id = $1
  ORDER BY created_at DESC;
```

**Delete follow:**
```sql
DELETE FROM feed_follows
  WHERE id = $1;
```

---

## Optimization Tips

### Query Performance

**1. Use EXPLAIN ANALYZE**
```sql
EXPLAIN ANALYZE
SELECT p.* FROM posts p
  JOIN feed_follows uf ON p.feed_id = uf.feed_id
  WHERE uf.user_id = '...'
  ORDER BY p.published_at DESC
  LIMIT 10;
```

**2. Add Missing Indexes**
```sql
-- If sequential scan appears in plan
CREATE INDEX idx_feed_follows_user_id ON feed_follows(user_id);
```

**3. Use LIMIT with pagination**
```sql
-- Avoid OFFSET on large result sets
SELECT * FROM posts
  WHERE feed_id = $1 AND id < $2
  ORDER BY published_at DESC
  LIMIT 20;
```

### Connection Pooling

**Configuration in app:**
```go
// Set connection pool size
sqlDB.SetMaxOpenConns(25)
sqlDB.SetMaxIdleConns(5)
sqlDB.SetConnMaxLifetime(5 * time.Minute)
```

### Batch Operations

**Instead of:**
```sql
INSERT INTO posts (...) VALUES (...);
INSERT INTO posts (...) VALUES (...);
INSERT INTO posts (...) VALUES (...);
```

**Use:**
```sql
INSERT INTO posts (...) VALUES 
  (...),
  (...),
  (...);
```

### Compression

**Reduce table size:**
```sql
-- Remove old posts
DELETE FROM posts
  WHERE published_at < NOW() - INTERVAL '30 days';

-- Vacuum to reclaim space
VACUUM ANALYZE posts;
```

### Monitoring Query Performance

```sql
-- Slow query log
SET log_min_duration_statement = 1000;  -- 1 second

-- Check current connections
SELECT datname, usename, application_name, state
  FROM pg_stat_activity;

-- Find long-running queries
SELECT query, query_start
  FROM pg_stat_statements
  WHERE query_start < NOW() - INTERVAL '1 minute'
  ORDER BY mean_exec_time DESC;
```

---

## Backup & Recovery

### Manual Backup

```bash
# Backup entire database
pg_dump $DATABASE_URL > backup.sql

# Backup specific table
pg_dump $DATABASE_URL -t users > users_backup.sql

# Backup to compressed file
pg_dump $DATABASE_URL | gzip > backup.sql.gz
```

### Restore from Backup

```bash
# Restore entire database
psql $DATABASE_URL < backup.sql

# Restore from compressed
gunzip -c backup.sql.gz | psql $DATABASE_URL
```

### Docker Backup

```bash
# Backup from container
docker-compose exec postgres pg_dump \
  -U postgres -d rss_generator > backup.sql

# Restore from container
docker-compose exec -T postgres psql \
  -U postgres -d rss_generator < backup.sql
```

### Automated Backups

**Backup script:**
```bash
#!/bin/bash
BACKUP_DIR="/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
FILE="$BACKUP_DIR/rss_generator_$TIMESTAMP.sql.gz"

pg_dump $DATABASE_URL | gzip > "$FILE"

# Keep only last 7 days
find $BACKUP_DIR -name "rss_generator_*.sql.gz" -mtime +7 -delete
```

Add to cron:
```bash
0 2 * * * /path/to/backup.sh  # Daily at 2 AM
```

---

## Migration Guide

### Adding New Column

1. **Create migration:**
```bash
goose -dir ./sql/migrations postgres "$DATABASE_URL" create add_column_to_users sql
```

2. **Write migration:**
```sql
-- +goose up
ALTER TABLE users ADD COLUMN email TEXT;
ALTER TABLE users ADD CONSTRAINT email_unique UNIQUE (email);

-- +goose down
ALTER TABLE users DROP COLUMN email;
```

3. **Apply:**
```bash
goose -dir ./sql/migrations postgres "$DATABASE_URL" up
```

### Adding New Table

1. **Create migration:**
```bash
goose -dir ./sql/migrations postgres "$DATABASE_URL" create create_comments_table sql
```

2. **Write migration:**
```sql
-- +goose up
CREATE TABLE comments (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    content TEXT NOT NULL,
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE
);

CREATE INDEX idx_comments_post_id ON comments(post_id);

-- +goose down
DROP TABLE comments;
```

3. **Update SQL queries:**
```bash
# sql/queries/comments.sql
-- name: CreateComment :one
INSERT INTO comments (id, content, post_id)
VALUES ($1, $2, $3)
RETURNING *;
```

4. **Generate code:**
```bash
sqlc generate
```

### Rolling Back Migrations

```bash
# Rollback last migration
goose -dir ./sql/migrations postgres "$DATABASE_URL" down

# Rollback to specific version
goose -dir ./sql/migrations postgres "$DATABASE_URL" down-to 001
```

---

## Database Maintenance

### Regular Maintenance Tasks

**Weekly:**
```sql
-- Vacuum and analyze
VACUUM ANALYZE;

-- Check table sizes
SELECT tablename, pg_size_pretty(pg_total_relation_size(tablename))
  FROM pg_tables
  WHERE schemaname = 'public'
  ORDER BY pg_total_relation_size(tablename) DESC;
```

**Monthly:**
```sql
-- Full vacuum (requires exclusive lock)
VACUUM FULL ANALYZE;

-- Reindex
REINDEX DATABASE rss_generator;

-- Update statistics
ANALYZE;
```

### Health Checks

```sql
-- Check for unused indexes
SELECT schemaname, tablename, indexname
  FROM pg_stat_user_indexes
  WHERE idx_scan = 0
  ORDER BY pg_relation_size(indexrelid) DESC;

-- Check bloat
SELECT schemaname, tablename, 
  pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
  FROM pg_tables
  WHERE schemaname = 'public'
  ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- Check constraints
SELECT constraint_name, constraint_type 
  FROM information_schema.table_constraints
  WHERE table_schema = 'public' AND table_name = 'users';
```

---

**Last Updated:** April 26, 2025
