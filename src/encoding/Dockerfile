FROM golang:1.24.5

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY src/go.mod src/go.sum ./
RUN go mod download

# Copy source code
COPY src/encoding encoding/.
COPY src/system_protocol system_protocol/.

# Build the Go app
RUN go build -o own_wiki encoding/main.go

# Run the app
CMD ["./own_wiki"]
