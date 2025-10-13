# Flux

```
███████╗██╗     ██╗   ██╗██╗  ██╗
██╔════╝██║     ██║   ██║╚██╗██╔╝
█████╗  ██║     ██║   ██║ ╚███╔╝ 
██╔══╝  ██║     ██║   ██║ ██╔██╗ 
██║     ███████╗╚██████╔╝██╔╝ ██╗
╚═╝     ╚══════╝ ╚═════╝ ╚═╝  ╚═╝
```

> **Flow through your tasks with ease.**

A modern task management API built with Go, Gin, and PostgreSQL.

## Prerequisites

- Go 1.22 or later
- Docker and Docker Compose

## Getting Started

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd flux
   ```

2. Start the services:
   ```bash
   docker-compose up --build
   ```

3. The API will be available at `http://localhost:8080`

## API Endpoints

### Health Check
- `GET /health` - Health check endpoint

### Users
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/:id` - Get a specific user
- `POST /api/v1/users` - Create a new user
- `PUT /api/v1/users/:id` - Update a user
- `DELETE /api/v1/users/:id` - Delete a user

### Tasks
- `GET /api/v1/tasks` - Get all tasks
- `GET /api/v1/tasks/:id` - Get a specific task
- `POST /api/v1/tasks` - Create a new task
- `PUT /api/v1/tasks/:id` - Update a task
- `DELETE /api/v1/tasks/:id` - Delete a task
- `GET /api/v1/users/:id/tasks` - Get all tasks for a specific user

## Development

### Running locally without Docker

1. Start the database:
   ```bash
   docker-compose up -d db
   ```

2. Run the application:
   ```bash
   go run main.go
   ```

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=taskdb
DB_PORT=5432
PORT=8080
```

## Project Structure

```
.
├── database/
│   ├── database.go    # Database connection
│   └── migrate.go     # Database migrations
├── handlers/
│   ├── task.go        # Task handlers
│   └── user.go        # User handlers
├── models/
│   ├── task.go        # Task model
│   └── user.go        # User model
├── routes/
│   └── routes.go      # API routes
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## Example API Usage

### Create a User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
```

### Create a Task
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Sample Task","description":"This is a test task","status":"pending","user_id":1}'
```

### Get All Tasks
```bash
curl http://localhost:8080/api/v1/tasks
```
