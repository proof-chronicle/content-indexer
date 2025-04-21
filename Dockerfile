# Stage 1: Build binary
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git make bash tzdata \
    # for headless Chrome (chromedp)
    chromium-chromedriver chromium \
    && rm -rf /var/cache/apk/*

WORKDIR /app

# Copy go.mod and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the indexer binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o indexer main.go

# Stage 2: Final image
FROM alpine:latest
RUN apk add --no-cache ca-certificates chromium-chromedriver chromium

WORKDIR /app
COPY --from=builder /app/indexer ./ 
COPY --from=builder /usr/lib/chromium /usr/lib/chromium

# Environment for Chrome
ENV CHROME_PATH=/usr/bin/chromium-browser \
    PATH="/usr/lib/chromium:${PATH}"

ENTRYPOINT ["/app/indexer"]