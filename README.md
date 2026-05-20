# GopherSocial

A production-ready social media application backend built with Go, featuring REST APIs for users, posts, and comments. This project demonstrates best practices in Go application architecture, database management, and containerization.

## Project Overview

GopherSocial is a backend service for a social media platform. It provides APIs for:

- User management and authentication
- Creating, reading, and managing posts
- Adding and managing comments on posts
- RESTful endpoints for all operations

## Technology Stack

- **Language**: Go 1.25.0
- **Web Framework**: [Chi Router](https://github.com/go-chi/chi)
- **Database**: PostgreSQL
- **Database Drivers**: pq (lib/pq)
- **Validation**: go-playground/validator
- **Containerization**: Docker & Docker Compose
- **Migrations**: golang-migrate

## Project Features

- **Optimistic Concurrency Control (OCC)**: Posts include versioning to prevent race conditions. When updating a post, the system checks if the version matches. If another client has modified the post (version mismatch), the update is rejected with a conflict error, ensuring data consistency without locking.
- **Clean Architecture**: Separation of concerns with cmd, internal, and store layers
- **Error Handling**: Custom error responses with meaningful messages
- **JSON Serialization**: Utilities for consistent JSON encoding/decoding
- **Input Validation**: Request validation using go-playground/validator
- **Database Abstraction**: Storage interface for flexible data persistence
- **Environment Configuration**: Flexible config via environment variables
- **Docker Support**: Containerized deployment ready
- **Database Migrations**: Version-controlled schema management

### Running the project with Docker Compose

Start the entire stack with PostgreSQL:

```bash
docker compose -f docker-compose.yaml up -d --build
```

The API will be available at `http://localhost:8080`

### Running Locally

1. **Start PostgreSQL** (ensure it's running on localhost:5432)

2. **Run database migrations**

```bash
make migrate-up
```

3. **Build and run the API**

```bash
go run ./cmd/api/main.go
```

The API will listen on `:8080` (default) or the port specified in `ADDR` environment variable.

## Database Migrations

Migrations are located in `cmd/migrate/migrations/` and managed using golang-migrate.

### Available Commands

```bash
# Run pending migrations
make migrate-up

# Rollback one migration
make migrate-down

# Create a new migration
make migrate-create migration_name
```
