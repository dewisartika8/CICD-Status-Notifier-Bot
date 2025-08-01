FROM golang:1.23-alpine

WORKDIR /app

# Copy go files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the simple version
RUN go build -o main ./cmd/simple_main.go

EXPOSE 8080

CMD ["./main"]