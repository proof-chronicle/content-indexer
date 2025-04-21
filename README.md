# Content Indexer Worker

A Go-based worker service responsible for:

- Listening to a RabbitMQ queue for new indexing tasks
- Rendering web pages (including those driven by JavaScript) via a headless browser
- Calculating a content hash (e.g. SHA-256) of the fully rendered page
- Sending a gRPC request to the Chain Gateway to anchor the hash on-chain
- Updating the corresponding database record with the resulting transaction ID and status

---

## Table of Contents

- [Content Indexer Worker](#content-indexer-worker)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Architecture](#architecture)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Usage](#usage)
  - [Docker](#docker)
  - [Development](#development)
  - [License](#license)

---

## Features

- **Queue-driven**: consumes tasks from RabbitMQ
- **Headless Browser**: uses Chrome/Chromium (via [chromedp](https://github.com/chromedp/chromedp)) to fully render pages
- **Content Hashing**: computes SHA-256 hash of the rendered HTML
- **gRPC Integration**: communicates with the Chain Gateway for on-chain anchoring
- **Database Update**: writes back transaction IDs and statuses to MySQL

---

## Architecture

```plaintext
[ RabbitMQ ] → [ content-indexer ] → [ Chain Gateway (gRPC) ]
                             ↘
                              [ MySQL ]
```

- **content-indexer**:
  - **Consumer**: subscribes to a RabbitMQ queue (`index_tasks`)
  - **Renderer**: launches headless Chrome to load and render JS-heavy pages
  - **Hasher**: extracts HTML and computes SHA-256
  - **Client**: dials Chain Gateway over gRPC to anchor the hash
  - **Updater**: executes SQL update on the `articles` table

---

## Prerequisites

- Go 1.20+
- Docker (for local development)
- RabbitMQ instance
- MySQL (or compatible) database
- Chain Gateway running and reachable via gRPC

---

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-org/content-indexer.git
   cd content-indexer
   ```
2. Fetch dependencies:
   ```bash
   go mod download
   ```

---

## Configuration

Copy the example environment file and adjust values:

```bash
cp .env.example .env
```

**.env**

```dotenv
# RabbitMQ
RABBITMQ_URL=amqp://rabbit:password@rabbitmq:5672/
QUEUE_NAME=index_tasks

# MySQL
DB_DSN=user:password@tcp(db:3306)/trustnews?parseTime=true

# Chain Gateway (gRPC)
CHAIN_GATEWAY_ADDR=chain-gateway:50051

# Headless Browser
BROWSER_TIMEOUT=30s
```

---

## Usage

```bash
# Run locally
go run ./cmd/indexer \
  --env .env \
  --log-level info

# Build and run binary
go build -o indexer ./cmd/indexer
./indexer --env .env
```

---

## Docker

Build and run with Docker Compose (from infra folder):

```yaml
# infra/docker-compose.yml
services:
  indexer:
    build:
      context: ../content-indexer
    env_file:
      - ../content-indexer/.env
    depends_on:
      - rabbitmq
      - db
```

Then:
```bash
cd infra
docker-compose up --build indexer
```

---

## Development

- Use `make lint` and `make test` to verify code quality.
- Ensure compatibility with latest Go version.
- Follow the Go project layout conventions.

---

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

