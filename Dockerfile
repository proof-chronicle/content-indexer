FROM golang:1.24

RUN go install github.com/air-verse/air@latest

WORKDIR /app

# Copy go.mod and go.sum first (for better cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of the app
COPY . .

EXPOSE 8080

CMD ["air"]
