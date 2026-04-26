# API Reference Documentation

Complete reference guide for all RSS Generator Project API endpoints with
detailed examples, parameters, and responses.

## Table of Contents

- [Base URL & Authentication](#base-url--authentication)
- [Request/Response Format](#requestresponse-format)
- [Status Codes](#status-codes)
- [Endpoints](#endpoints)
  - [Health & Status](#health--status)
  - [User Management](#user-management)
  - [Feed Management](#feed-management)
  - [Feed Subscriptions](#feed-subscriptions)
  - [Posts](#posts)
- [Error Responses](#error-responses)
- [Rate Limiting](#rate-limiting)
- [Examples](#examples)

---

## Base URL & Authentication

### Base URL

```
http://localhost:8080/v1
```

For production, replace `localhost:8080` with your domain and port.

### Authentication

The API uses API key-based authentication. Include your API key in the
`Authorization` header:

```
Authorization: ApiKey {api_key}
```

**Example:**

```bash
curl -H "Authorization: ApiKey a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6" \
  http://localhost:8080/v1/user
```

### Obtaining an API Key

1. Create a user account via `POST /user`
2. The response will include your unique `api_key`
3. Store this key securely
4. Include it in all authenticated requests

---

## Request/Response Format

### Content-Type

All endpoints accept and return JSON:

```
Content-Type: application/json
```

### Request Body

For endpoints accepting data, send JSON in the request body:

```bash
curl -X POST http://localhost:8080/v1/user \
  -H "Content-Type: application/json" \
  -d '{"name": "john_doe"}'
```

### Response Format

All responses follow this structure:

**Success Response:**

```json
{
  "id": "uuid",
  "field1": "value1",
  "field2": "value2"
}
```

**Error Response:**

```json
{
  "error": "Error description"
}
```

---

## Status Codes

| Code                        | Meaning               | Usage                                           |
| --------------------------- | --------------------- | ----------------------------------------------- |
| `200 OK`                    | Success               | GET requests, successful updates                |
| `201 Created`               | Resource created      | POST requests that create resources             |
| `204 No Content`            | Success with no body  | DELETE requests, batch operations               |
| `400 Bad Request`           | Invalid input         | Malformed JSON, validation errors, logic errors |
| `403 Forbidden`             | Authentication failed | Missing/invalid API key, malformed auth header  |
| `404 Not Found`             | Resource not found    | Requesting non-existent resource                |
| `500 Internal Server Error` | Server error          | Unexpected server error, database issues        |

---

## Endpoints

### Health & Status

#### Check API Health

Verify that the API server is running and healthy.

```http
GET /healthz
```

**Parameters:** None

**Response:**

```json
{}
```

**Status Code:** `200 OK`

**Example:**

```bash
curl http://localhost:8080/v1/healthz
```

**Response:**

```
{}
```

---

### User Management

#### Create User

Register a new user account.

```http
POST /user
Content-Type: application/json

{
  "name": "username"
}
```

**Request Parameters:**

| Field  | Type   | Required | Description                                           |
| ------ | ------ | -------- | ----------------------------------------------------- |
| `name` | string | Yes      | Unique username (alphanumeric, no spaces recommended) |

**Response:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-04-26T10:30:00Z",
  "updated_at": "2025-04-26T10:30:00Z",
  "name": "john_doe",
  "api_key": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
}
```

**Status Code:** `201 Created`

**Error Responses:**

```json
{
  "error": "Error Parsing json: EOF"
}
```

Status: `400 Bad Request` - Invalid JSON

```json
{
  "error": "Couldn't create user: duplicate username"
}
```

Status: `400 Bad Request` - Username already exists

**Example:**

```bash
curl -X POST http://localhost:8080/v1/user \
  -H "Content-Type: application/json" \
  -d '{"name": "alice"}'
```

**Response:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-04-26T10:30:00Z",
  "updated_at": "2025-04-26T10:30:00Z",
  "name": "alice",
  "api_key": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
}
```

---

#### Get Current User

Retrieve the authenticated user's profile information.

```http
GET /user
Authorization: ApiKey {api_key}
```

**Parameters:** None (authentication via header)

**Response:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-04-26T10:30:00Z",
  "updated_at": "2025-04-26T10:30:00Z",
  "name": "alice",
  "api_key": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
}
```

**Status Code:** `200 OK`

**Error Responses:**

```json
{
  "error": "Auth error: No Authentication info found"
}
```

Status: `403 Forbidden` - Missing Authorization header

```json
{
  "error": "Auth error: Malformed auth Header"
}
```

Status: `403 Forbidden` - Invalid header format

```json
{
  "error": "Couldn't get user {error details}"
}
```

Status: `400 Bad Request` - Database error

**Example:**

```bash
curl -H "Authorization: ApiKey a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6" \
  http://localhost:8080/v1/user
```

---

#### Delete User

Deactivate and delete a user account and all associated data.

```http
DELETE /user
Authorization: ApiKey {api_key}
```

**Parameters:** None (authentication via header)

**Response:**

```json
{}
```

**Status Code:** `204 No Content`

**Error Responses:**

```json
{
  "error": "Auth error: No Authentication info found"
}
```

Status: `403 Forbidden` - Missing authentication

**Cascading Deletes:**

- User's feeds
- User's feed follows
- Posts from user's feeds (if no other followers)

**Example:**

```bash
curl -X DELETE \
  -H "Authorization: ApiKey a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6" \
  http://localhost:8080/v1/user
```

**Warning:** This operation is permanent and cannot be undone. All user data
will be deleted.

---

### Feed Management

#### Create Feed

Add a new RSS feed source to the system.

```http
POST /feed
Authorization: ApiKey {api_key}
Content-Type: application/json

{
  "name": "Feed Name",
  "url": "https://example.com/feed/rss"
}
```

**Request Parameters:**

| Field  | Type   | Required | Description                   |
| ------ | ------ | -------- | ----------------------------- |
| `name` | string | Yes      | Descriptive name for the feed |
| `url`  | string | Yes      | Valid RSS feed URL            |

**Response:**

```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-04-26T10:35:00Z",
  "updated_at": "2025-04-26T10:35:00Z",
  "name": "Tech News Daily",
  "url": "https://example.com/feed/rss",
  "user_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Status Code:** `201 Created`

**Error Responses:**

```json
{
  "error": "Auth error: No Authentication info found"
}
```

Status: `403 Forbidden` - Missing authentication

```json
{
  "error": "Error Parsing json: missing required field"
}
```

Status: `400 Bad Request` - Invalid JSON

```json
{
  "error": "Couldn't create user: {error details}"
}
```

Status: `400 Bad Request` - Database error

**Example:**

```bash
curl -X POST \
  -H "Authorization: ApiKey a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Hacker News",
    "url": "https://news.ycombinator.com/rss"
  }' \
  http://localhost:8080/v1/feed
```

---

#### Get All Feeds

Retrieve all available feeds in the system.

```http
GET /feed
```

**Parameters:** None (public endpoint)

**Response:**

```json
[
  {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "created_at": "2025-04-26T10:35:00Z",
    "updated_at": "2025-04-26T10:35:00Z",
    "name": "Tech News Daily",
    "url": "https://example.com/feed/rss",
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
  },
  {
    "id": "770e8400-e29b-41d4-a716-446655440001",
    "created_at": "2025-04-26T10:40:00Z",
    "updated_at": "2025-04-26T10:40:00Z",
    "name": "Dev Blog",
    "url": "https://devblog.example.com/feed",
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
  }
]
```

**Status Code:** `200 OK`

**Empty Response:**

```json
[]
```

Returns an empty array if no feeds exist.

**Example:**

```bash
curl http://localhost:8080/v1/feed
```

---

### Feed Subscriptions

#### Create Feed Follow

Subscribe to a specific feed.

```http
POST /feed-follow
Authorization: ApiKey {api_key}
Content-Type: application/json

{
  "feed_id": "660e8400-e29b-41d4-a716-446655440000"
}
```

**Request Parameters:**

| Field     | Type | Required | Description              |
| --------- | ---- | -------- | ------------------------ |
| `feed_id` | UUID | Yes      | ID of the feed to follow |

**Response:**

```json
{
  "id": "880e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-04-26T10:45:00Z",
  "updated_at": "2025-04-26T10:45:00Z",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "feed_id": "660e8400-e29b-41d4-a716-446655440000"
}
```

**Status Code:** `201 Created`

**Error Responses:**

```json
{
  "error": "Auth error: No Authentication info found"
}
```

Status: `403 Forbidden` - Missing authentication

```json
{
  "error": "Error Parsing Json"
}
```

Status: `500 Internal Server Error` - Invalid JSON

```json
{
  "error": "Error Creating a Feed Follow"
}
```

Status: `500 Internal Server Error` - Database error (e.g., user already
following)

**Example:**

```bash
curl -X POST \
  -H "Authorization: ApiKey a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6" \
  -H "Content-Type: application/json" \
  -d '{
    "feed_id": "660e8400-e29b-41d4-a716-446655440000"
  }' \
  http://localhost:8080/v1/feed-follow
```

---

#### Get Feed Follows

Retrieve all feeds a user is following.

```http
GET /feed-follow
Authorization: ApiKey {api_key}
```

**Parameters:** None (authentication via header)

**Response:**

```json
[
  {
    "id": "880e8400-e29b-41d4-a716-446655440000",
    "created_at": "2025-04-26T10:45:00Z",
    "updated_at": "2025-04-26T10:45:00Z",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "feed_id": "660e8400-e29b-41d4-a716-446655440000"
  },
  {
    "id": "990e8400-e29b-41d4-a716-446655440001",
    "created_at": "2025-04-26T10:50:00Z",
    "updated_at": "2025-04-26T10:50:00Z",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "feed_id": "770e8400-e29b-41d4-a716-446655440001"
  }
]
```

**Status Code:** `200 OK`

**Empty Response:**

```json
[]
```

Returns an empty array if user is not following any feeds.

**Example:**

```bash
curl -H "Authorization: ApiKey a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6" \
  http://localhost:8080/v1/feed-follow
```

---

#### Delete Feed Follow

Unsubscribe from a feed.

```http
DELETE /feed-follow/{feed-follow-id}
Authorization: ApiKey {api_key}
```

**Path Parameters:**

| Parameter        | Type | Description                                  |
| ---------------- | ---- | -------------------------------------------- |
| `feed-follow-id` | UUID | ID of the feed follow relationship to delete |

**Response:**

```json
{}
```

**Status Code:** `204 No Content`

**Error Responses:**

```json
{
  "error": "Auth error: No Authentication info found"
}
```

Status: `403 Forbidden` - Missing authentication

```json
{
  "error": "Error deleting feed follow"
}
```

Status: `400 Bad Request` - Feed follow not found or doesn't belong to user

**Example:**

```bash
curl -X DELETE \
  -H "Authorization: ApiKey a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6" \
  http://localhost:8080/v1/feed-follow/880e8400-e29b-41d4-a716-446655440000
```

---

### Posts

#### Get User Posts

Retrieve posts from all feeds the user is following (paginated).

```http
GET /posts
Authorization: ApiKey {api_key}
```

**Parameters:** None (currently limited to 10 posts, no pagination control)

**Response:**

```json
[
  {
    "id": "aa0e8400-e29b-41d4-a716-446655440000",
    "created_at": "2025-04-26T11:00:00Z",
    "updated_at": "2025-04-26T11:00:00Z",
    "title": "Breaking: New Tech Release",
    "description": "A groundbreaking technology was just released...",
    "published_at": "2025-04-26T10:50:00Z",
    "url": "https://example.com/article/123",
    "feed_id": "660e8400-e29b-41d4-a716-446655440000"
  },
  {
    "id": "bb0e8400-e29b-41d4-a716-446655440001",
    "created_at": "2025-04-26T11:05:00Z",
    "updated_at": "2025-04-26T11:05:00Z",
    "title": "Industry Leaders Meet",
    "description": null,
    "published_at": "2025-04-26T09:30:00Z",
    "url": "https://example.com/article/456",
    "feed_id": "770e8400-e29b-41d4-a716-446655440001"
  }
]
```

**Status Code:** `200 OK`

**Empty Response:**

```json
[]
```

Returns an empty array if:

- User is not following any feeds
- No posts have been scraped yet
- User hasn't reached the posts they're looking for

**Fields Description:**

| Field          | Type              | Description                        |
| -------------- | ----------------- | ---------------------------------- |
| `id`           | UUID              | Unique post identifier             |
| `created_at`   | string (ISO 8601) | When post was added to database    |
| `updated_at`   | string (ISO 8601) | Last modification time             |
| `title`        | string            | Post title from RSS feed           |
| `description`  | string \| null    | Post content/summary (may be null) |
| `published_at` | string (ISO 8601) | Publication date from RSS feed     |
| `url`          | string            | URL to the full article            |
| `feed_id`      | UUID              | Which feed this post came from     |

**Error Responses:**

```json
{
  "error": "Auth error: No Authentication info found"
}
```

Status: `403 Forbidden` - Missing authentication

```json
{
  "error": "Couldn't get posts for user: {error details}"
}
```

Status: `400 Bad Request` - Database error

**Example:**

```bash
curl -H "Authorization: ApiKey a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6" \
  http://localhost:8080/v1/posts
```

---

## Error Responses

### Common Error Scenarios

#### Missing or Invalid API Key

```http
GET /user
Authorization: ApiKey invalid-key
```

```json
{
  "error": "Couldn't get user: user not found"
}
```

Status: `400 Bad Request`

---

#### Malformed Authorization Header

```http
GET /user
Authorization: Bearer token-here
```

```json
{
  "error": "Auth error: Malformed auth Header"
}
```

Status: `403 Forbidden`

**Valid format:** `Authorization: ApiKey {key}`

---

#### Missing Authorization Header

```http
GET /user
```

```json
{
  "error": "Auth error: No Authentication info found"
}
```

Status: `403 Forbidden`

---

#### Duplicate Username

```http
POST /user
Content-Type: application/json

{"name": "alice"}
```

(if alice already exists)

```json
{
  "error": "Couldn't create user: duplicate username"
}
```

Status: `400 Bad Request`

---

#### Duplicate Feed Follow

```http
POST /feed-follow
Authorization: ApiKey {api_key}
Content-Type: application/json

{"feed_id": "660e8400-e29b-41d4-a716-446655440000"}
```

(if already following this feed)

```json
{
  "error": "Error Creating a Feed Follow"
}
```

Status: `500 Internal Server Error`

---

## Rate Limiting

**Current Rate Limiting:** None implemented

**Recommendations for Production:**

- Implement per-user rate limiting
- Suggested: 100 requests per minute per API key
- 1000 requests per minute per IP address
- Return `429 Too Many Requests` when limit exceeded

---

## Examples

### Complete User Journey

#### 1. Create a User

```bash
curl -X POST http://localhost:8080/v1/user \
  -H "Content-Type: application/json" \
  -d '{"name": "alice"}'
```

Response:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-04-26T10:30:00Z",
  "updated_at": "2025-04-26T10:30:00Z",
  "name": "alice",
  "api_key": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
}
```

Save the API key: `a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6`

---

#### 2. Get All Feeds

```bash
curl http://localhost:8080/v1/feed
```

Response:

```json
[
  {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "created_at": "2025-04-26T10:35:00Z",
    "updated_at": "2025-04-26T10:35:00Z",
    "name": "Hacker News",
    "url": "https://news.ycombinator.com/rss",
    "user_id": "550e8400-e29b-41d4-a716-446655440001"
  }
]
```

---

#### 3. Follow a Feed

```bash
curl -X POST http://localhost:8080/v1/feed-follow \
  -H "Authorization: ApiKey a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6" \
  -H "Content-Type: application/json" \
  -d '{"feed_id": "660e8400-e29b-41d4-a716-446655440000"}'
```

Response:

```json
{
  "id": "880e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-04-26T10:45:00Z",
  "updated_at": "2025-04-26T10:45:00Z",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "feed_id": "660e8400-e29b-41d4-a716-446655440000"
}
```

---

#### 4. View Posts From Followed Feeds

```bash
curl -H "Authorization: ApiKey a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6" \
  http://localhost:8080/v1/posts
```

Response:

```json
[
  {
    "id": "aa0e8400-e29b-41d4-a716-446655440000",
    "created_at": "2025-04-26T11:00:00Z",
    "updated_at": "2025-04-26T11:00:00Z",
    "title": "Breaking News Story",
    "description": "Details about the news...",
    "published_at": "2025-04-26T10:50:00Z",
    "url": "https://example.com/article/123",
    "feed_id": "660e8400-e29b-41d4-a716-446655440000"
  }
]
```

---

#### 5. View Your Profile

```bash
curl -H "Authorization: ApiKey a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6" \
  http://localhost:8080/v1/user
```

Response:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-04-26T10:30:00Z",
  "updated_at": "2025-04-26T10:30:00Z",
  "name": "alice",
  "api_key": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
}
```

---

## Client Library Examples

### JavaScript/TypeScript

```javascript
class RSSGeneratorClient {
  constructor(baseUrl, apiKey) {
    this.baseUrl = baseUrl;
    this.apiKey = apiKey;
  }

  async request(method, endpoint, body = null) {
    const headers = {
      "Content-Type": "application/json",
      Authorization: `ApiKey ${this.apiKey}`,
    };

    const options = { method, headers };
    if (body) options.body = JSON.stringify(body);

    const response = await fetch(`${this.baseUrl}${endpoint}`, options);
    return response.json();
  }

  async getUser() {
    return this.request("GET", "/user");
  }

  async getFeeds() {
    return fetch(`${this.baseUrl}/feed`).then((r) => r.json());
  }

  async followFeed(feedId) {
    return this.request("POST", "/feed-follow", { feed_id: feedId });
  }

  async getPosts() {
    return this.request("GET", "/posts");
  }
}

// Usage
const client = new RSSGeneratorClient(
  "http://localhost:8080/v1",
  "api_key_here",
);
const user = await client.getUser();
const posts = await client.getPosts();
```

### Python

```python
import requests

class RSSGeneratorClient:
    def __init__(self, base_url, api_key):
        self.base_url = base_url
        self.api_key = api_key
        self.headers = {
            'Content-Type': 'application/json',
            'Authorization': f'ApiKey {api_key}'
        }

    def get_user(self):
        return requests.get(f'{self.base_url}/user', headers=self.headers).json()

    def get_feeds(self):
        return requests.get(f'{self.base_url}/feed').json()

    def follow_feed(self, feed_id):
        return requests.post(
            f'{self.base_url}/feed-follow',
            json={'feed_id': feed_id},
            headers=self.headers
        ).json()

    def get_posts(self):
        return requests.get(f'{self.base_url}/posts', headers=self.headers).json()

# Usage
client = RSSGeneratorClient('http://localhost:8080/v1', 'api_key_here')
user = client.get_user()
posts = client.get_posts()
```

---

**Last Updated:** April 26, 2025
