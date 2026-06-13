# GopherSocial

Social media backend built with Go. Features REST APIs for users, posts, comments, and followers — with pagination, search, filtering, and Swagger documentation.

---

## Screenshots

### Swagger UI

![Swagger UI]()

---

## Features

- Email - Email is sent for account activation
- Posts — create, read, update, delete
- Comments — add and fetch comments on posts
- Users — fetch user profiles
- Followers — follow and unfollow users, get followers and following lists
- Feed — paginated user feed with search and filter
- Optimistic concurrency control on post updates
- Full text search on posts using trigram indexes
- Swagger API documentation
- Database migrations
- Hot reload with Air
- Containerized with Docker
- We will be adding Redis for user profile caching.

---

## Technology Stack

| Tool                    | Purpose                         |
| ----------------------- | ------------------------------- |
| Go 1.25                 | Language                        |
| Chi Router              | HTTP routing and middleware     |
| PostgreSQL 16           | Database                        |
| lib/pq                  | Postgres driver                 |
| golang-migrate          | Database migrations             |
| go-playground/validator | Request validation              |
| swaggo/http-swagger     | Swagger UI and docs generation  |
| Air                     | Hot reload during development   |
| Docker + Docker Compose | Containerization                |
| direnv                  | Environment variable management |

## Running With Docker

```zsh
# clone the repo
git clone <your-repo-url>
cd GopherSocial

# start everything — postgres + backend
docker-compose up --build

# API available at
http://localhost:8080

# Swagger docs at
http://localhost:8080/v1/swagger/index.html
```

---

## Set Up Environment

Your `.envrc` should contain:

```zsh
export ADDR="your_sweet_port"
export EXTERNAL_URL="localhost:your_sweet_port"
export DB_ADDR="postgres://username:password@localhost:5432/databasename?sslmode=disable"
export DB_MAX_OPEN_CONNS=30
export DB_MAX_IDLE_CONNS=30
export DB_MAX_IDLE_TIME="15m"
export Go_ENV="development"
```

## Database Migrations

Migrations live in `cmd/migrate/migrations/` and are managed with golang-migrate.

```zsh
# run all pending migrations
make migrate-up

# rollback one migration
make migrate-down

# create a new migration
make migrate-create migration_name

# force a version (when dirty database)
make migrate-force version=6
```

---

## Architecture Decisions

### Repository Pattern

All database operations go through a `Storage` interface. Handlers never touch the database directly — they call methods on the interface. This makes the code testable and allows swapping implementations without changing handlers.

```
handler → Storage interface → PostStore → postgres
```

### Optimistic Concurrency Control

Post updates use optimistic concurrency control to handle concurrent edits without locking rows.

Every post has a `version` column starting at `1`. When updating, the query checks `WHERE id=$1 AND version=$2`. If someone else updated first, version changes and the update hits `0` rows — a `409 Conflict` is returned.

```
User 1 reads post → version=1
User 2 reads post → version=1

User 1 updates → version matches → saved → version becomes 2
User 2 updates → version is now 2 → 0 rows affected → 409 Conflict
```

Why optimistic over pessimistic? Concurrent edits on the same post are rare in a social media app. Optimistic control avoids row locking which would slow reads under high traffic.

### Dependency Injection

The `application` struct receives its dependencies — `store`, `config` — from `main.go`. Nothing creates its own dependencies. This keeps the code loosely coupled and easy to test.

### Database Indexes

Indexes are added for all frequently searched and joined columns:

```sql
-- full text search on title and content
CREATE INDEX idx_posts_title ON posts USING gin (title gin_trgm_ops);

-- array search on tags
CREATE INDEX idx_posts_tags ON posts USING gin (tags);

-- exact match on foreign keys
CREATE INDEX idx_posts_user_id     ON posts    (user_id);
CREATE INDEX idx_comments_post_id  ON comments (post_id);
CREATE INDEX idx_users_username    ON users    (username);
```

---
