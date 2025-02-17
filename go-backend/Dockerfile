# Backend Dockerfile
FROM golang:1.24.0-alpine

WORKDIR /app

# Copy go mod files
COPY go.* ./

# Download dependencies
RUN go mod download

# Copy the entire backend directory
COPY . .

EXPOSE 3000 

# Use go run instead of building binary
CMD ["go", "run", "."]
