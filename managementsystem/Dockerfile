# Use official Golang image as a base
FROM golang:1.21-alpine

# Install Air for live reload
RUN go install github.com/cosmtrek/air@v1.42.0

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the application code
COPY . .

# Build the Go application
RUN go build -o /app/main ./cmd

# Verify that the binary was created
RUN ls -l /app

# Expose application port
EXPOSE 8080

# Start the application
CMD ["air"]

