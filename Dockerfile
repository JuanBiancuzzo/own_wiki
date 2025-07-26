FROM golang:1.24.5

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go app
RUN go build -o own_wiki crear_db.go

# Run the app
CMD ["./own_wiki"]
