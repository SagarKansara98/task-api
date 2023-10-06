**# Task Manager RESTful API using Golang, Echo, and PostgreSQL**

This repository contains a RESTful API implemented in Go using the Echo framework and PostgreSQL database for managing tasks. The API allows users to register, log in, create, read, update, and delete tasks. Users can only access their own tasks after authentication. Concurrency is demonstrated through Goroutines and channels for marking tasks as "done" concurrently.

## Table of Contents
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
- [Authentication & Authorization](#authentication--authorization)
- [Task Management](#task-management)
  - [Endpoints](#endpoints)
  - [Task Data](#task-data)
- [Database Integration](#database-integration)
- [Middleware & Error Handling](#middleware--error-handling)
- [Concurrency](#concurrency)
- [Documentation](#documentation)
- [Postman Collection](#postman-collection)

## Getting Started

### Prerequisites
- Go installed on your system. You can download it from [here](https://golang.org/dl/).

### Installation
1. Clone the repository:
   ```sh
   git clone git@github.com:SagarKansara98/task-api.git
   ```
2. Change directory to the project folder:
   ```sh
   cd task-api/app
   ```
3. Install dependencies:
   ```sh
   go mod tidy; go mod vendor
   ```
   
### Configuration
1. Create a PostgreSQL database and update the database configuration in `config/config.go`.
2. Set environment variables for authentication JWT secret and other sensitive data.
3. Run the migration (as time permit i have not configure migrate cli so it's just like run the query present in migration folder)

## Authentication & Authorization
- User registration and login endpoints are provided.
- Authentication is implemented using JWT (JSON Web Tokens).
- Users can only access their own tasks after authentication.

## Task Management

### Endpoints

- `POST /api/v1/user`: Register a new user.
- `POST /api/v1/user/login`: Log in with registered credentials.
- `GET /api/v1/user`: return login user details.
- `PUT /api/v1/user`: update login user details.
- `Delete /api/v1/user`: delete login user.
---
- `GET /api/v1/task`: Get all tasks for the authenticated user.
- `POST /api/v1/task`: Create a new task.
- `GET /api/v1/task/:id`: Get a specific task by ID.
- `PUT /api/v1/task/:id`: Update a task by ID.
- `DELETE /api/v1/task/:id`: Delete a task by ID.
- `PATCH /api/v1/task/mark-as-done`: Mark multiple tasks as "done" concurrently.

### Task Data
- Each task must have a `title`, `description`, and `status`.
- Valid `status` values: "todo", "in progress", "done".

## Database Integration
- PostgreSQL database is used for persistent storage.

## Middleware & Error Handling
- Middleware is implemented for authentication and logging incoming requests.
- Errors are handled gracefully, providing meaningful error responses with appropriate status codes.

## Concurrency
- The API allows marking multiple tasks as "done" concurrently using Goroutines and channels.
- Concurrency is implemented in the `POST /api/v1/mark-as-done` endpoint.

## Documentation
1. **How to Run the Application:**
   - Set up the PostgreSQL database and configure the database connection in `config/config.go`.
   - Set environment variables for sensitive data.
   - Run the application using the following command: (please add path to app directotry)
     ```sh
     go run main.go
     ```
3. **Authentication:**
   - User registration and login are required for accessing task-related endpoints.
   - Use the JWT token obtained after login for authentication. Include the token in the Authorization header as `Bearer <token>` for authenticated requests.
4. **Error Responses:**
   - API endpoints provide meaningful error responses with appropriate status codes and error messages.
5. **Concurrency:**
   - The `POST /api/v1/task/mark-as-done` endpoint allows marking multiple tasks as "done" concurrently using Goroutines and channels.

## Postman Collection
Import the provided [Postman Collection](https://documenter.getpostman.com/view/29493818/2s9YJezMkG) to test the API endpoints directly.