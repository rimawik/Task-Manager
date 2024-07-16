# Task Manager

Task Manager is a simple task management application built in Go. It allows users to manage their tasks with basic CRUD operations.

## Features

- **Create**: Add new tasks with a title and description.
- **Read**: View all tasks or a specific task by ID.
- **Update**: Modify existing tasks.
- **Delete**: Remove tasks from the list.

## Technologies Used

- **Go**: Programming language used for development.
- **PostgreSQL**: Database management system for storing tasks.
- **Redis**: In-memory data structure store for caching.
- **Gorilla Mux**: HTTP router used for handling routing and middleware.
- **Swagger**: API documentation tool.
- **Testing**: Unit tests written using the standard Go testing package.


## Getting Started

### Prerequisites

To run this project, you need to have the following software installed on your system:

- Go
- PostgreSQL
- Redis
- swag

### Installation

1. **Clone the repository**:

2. **Set up the database**:

Ensure PostgreSQL and Redis are running.
Update configuration in config/config.go.

3. **Build and Run**:

Build and run the application:
```bash
make build
./bin/main
```

4. **Generate Swagger**:

```bash
make swagger
```

4. **Access The API Documentation**:
Open your browser and go to http://localhost:8080/swagger/index.html
