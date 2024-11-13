# Real-Time Data Streaming API with Golang and Redpanda (Kafka)

## Table of Contents
1. [Project Overview](#project-overview)
2. [Features](#features)
3. [Prerequisites](#prerequisites)
4. [Project Structure](#project-structure)
5. [Setup and Installation](#setup-and-installation)
6. [Running the Application](#running-the-application)
7. [API Endpoints](#api-endpoints)
8. [Authentication](#authentication)
9. [Running Tests](#running-tests)
10. [Performance Benchmarking](#performance-benchmarking)
11. [Monitoring and Logging](#monitoring-and-logging)
12. [Rate Limiting](#rate-limiting)
13. [Error Handling](#error-handling)
14. [Scalability and Performance](#scalability-and-performance)

## Project Overview

This project implements a high-performance API in Golang that streams data to and from Redpanda (Kafka). The API allows clients to stream data in real-time and receive processed results back through a bi-directional communication channel. It's designed to handle 1000+ concurrent streams efficiently, making it suitable for high-throughput, low-latency data processing applications.

## Features

- RESTful API supporting real-time data streaming
- Kafka (Redpanda) integration for message processing
- Concurrent handling of multiple data streams
- Real-time data processing with WebSocket support
- API key authentication for secure access
- Rate limiting to prevent abuse and ensure fair usage
- Comprehensive error handling and structured logging
- Performance benchmarking tools for system evaluation

## Prerequisites

- Go 1.16 or higher
- Redpanda (or Kafka) cluster
- Git (for cloning the repository)


## Setup and Installation

1. Clone the repository: git clone https://github.com/aryanc1027/api-kafka-golang.git

2. Install Dependencies: go mod tidy


3. Set up Redpanda (or Kafka):
- Follow the official Redpanda documentation to set up a cluster
- Note down the broker addresses for configuration

4. Configure the application:
- Copy `config.example.yaml` to `config.yaml`
- Update the Kafka hosts and other settings in `config.yaml`

## Running the Application

1. Start the API server: go run cmd/api/main.go


2. The server will start on the configured port (default: 8080)

## API Endpoints

1. Start a new data stream
- Method: POST
- Endpoint: `/api/stream/start`
- Response: Returns a unique `stream_id`

2. Send data to a specific stream
- Method: POST
- Endpoint: `/api/stream/{stream_id}/send`
- Body: JSON data to be processed

3. Retrieve real-time results (WebSocket)
- Method: GET
- Endpoint: `/api/stream/{stream_id}/results`
- Description: Establishes a WebSocket connection for real-time results

Example usage:

```bash
# Start a new stream
curl -X POST http://localhost:8080/api/stream/start -H "X-API-Key: your-api-key"

# Send data to a stream
curl -X POST http://localhost:8080/api/stream/your-stream-id/send \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{"data": "your data here"}'

# For WebSocket connection, use a WebSocket client to connect to:
ws://localhost:8080/api/stream/your-stream-id/results
```
## Authentication

All API endpoints are protected with API key authentication. Include the API key in the X-API-Key header for all requests. The API key should be configured in the config.yaml file.
Running

## Performance Benchmarking 

1. Ensure API server is running
2. Run the benchmark script: go run scripts/benchmark.go

## Rate Limiting

The API implements rate limiting to prevent abuse. The current rate limits are:
- **1 request per second** per IP address
- **Burst allowance of 5 requests**

These limits can be adjusted in the `internal/api/middleware.go` file.

## Error Handling

The API provides detailed error messages and appropriate HTTP status codes for various error scenarios:

- **400 Bad Request**: Invalid input data
- **401 Unauthorized**: Invalid or missing API key
- **404 Not Found**: Stream not found
- **429 Too Many Requests**: Rate limit exceeded
- **500 Internal Server Error**: Unexpected server errors

All errors are logged for debugging and monitoring purposes.

## Scalability and Performance

The application is designed to efficiently handle **1000+ concurrent streams**. Key performance features include:

- **Goroutines** for concurrent stream processing
- **Kafka producers and consumers** for efficient message processing
- **WebSocket** for real-time, bi-directional communication
- **Buffered channels** for optimized message passing

To further scale the application:
1. Deploy multiple instances behind a **load balancer**.
2. Scale the **Kafka/Redpanda cluster**.
3. Use a distributed cache (e.g., **Redis**) for sharing state between instances.






