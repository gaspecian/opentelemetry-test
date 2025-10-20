# TODO API with OpenTelemetry and MongoDB

A simple TODO CRUD API built with Go, Gorilla Mux, MongoDB, and OpenTelemetry instrumentation.

## Setup

1. Start MongoDB with Docker Compose:
```bash
docker-compose up -d
```

2. Install Go dependencies:
```bash
go mod tidy
```

## Run the application

```bash
go run main.go
```

## API Endpoints

- `GET /todos` - Get all todos
- `POST /todos` - Create a new todo
- `GET /todos/{id}` - Get a specific todo (use MongoDB ObjectID)
- `PUT /todos/{id}` - Update a todo
- `DELETE /todos/{id}` - Delete a todo

## Example requests

```bash
# Create a todo
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn OpenTelemetry","done":false}'

# Get all todos
curl http://localhost:8080/todos

# Update a todo (replace {id} with actual ObjectID from create response)
curl -X PUT http://localhost:8080/todos/{id} \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn OpenTelemetry","done":true}'

# Delete a todo
curl -X DELETE http://localhost:8080/todos/{id}
```

## MongoDB Connection

- Host: localhost:27017
- Username: admin
- Password: password
- Database: todoapp
- Collection: todos

OpenTelemetry traces will be printed to stdout.
