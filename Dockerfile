FROM golang:1.24.5

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY src/go.mod src/go.sum ./
RUN go mod download

# Copy variables
COPY .env .env

# Copy source code
COPY src/ .

# Build the Go app
RUN go build -o own_wiki main.go

# Run the app
CMD ["./own_wiki"]
