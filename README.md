# Orders CRUD API

## Overview

The Orders CRUD API is a Go application that provides a RESTful interface for managing orders. It uses Redis as an in-memory database to store and retrieve order data efficiently. The application is built using the `go-chi/chi` router, which allows for clean and idiomatic routing patterns.

## Endpoints

The API exposes the following endpoints for managing orders:

- `POST /orders`: Create a new order.
- `GET /orders`: Retrieve a list of all orders.
- `GET /orders/{id}`: Get the details of an order by its ID.
- `PUT /orders/{id}`: Update an existing order by its ID.
- `DELETE /orders/{id}`: Delete an order by its ID.

## Technologies

- **Go**: The primary programming language used to build the application.
- **go-chi/chi**: A lightweight and composable router for building Go HTTP services.
- **Redis**: An in-memory data structure store used for caching and data persistence.

## Setup and Installation

1. Ensure you have Go installed on your machine.
2. Clone the repository to your local machine.
3. Navigate to the project directory and run `go mod download` to download dependencies.
4. Run the application with `go run main.go`.

## Usage

After starting the application, you can interact with the API using tools like `curl` or Postman, or by using the provided client libraries.

## Development

For development purposes, you can run Redis locally or as a Docker container. Instructions for both methods are provided above.
