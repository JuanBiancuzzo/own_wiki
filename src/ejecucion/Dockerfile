FROM golang:1.24.5

# Set working directory
WORKDIR /app

# Copy go mod and sum files
RUN go install github.com/air-verse/air@latest

# Copy variables
COPY .env .env

# Copy go mod and sum files
COPY src/go.mod src/go.sum ./
RUN go mod download


CMD ["air", "-c", ".air.toml"]