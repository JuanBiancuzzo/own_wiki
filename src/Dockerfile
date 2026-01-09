FROM golang:1.25.5-alphine3.23 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

ARG TARGET_PLATFORM = "terminal" 

# Genera los scripts de grpc
RUN protoc --proto_path=proto proto/*.proto --go_out=. --go-grpc_out=.

ARG TARGETOS
ARG TARGETARCH

# Compilamos el exe
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -tags $TARGET_PLATFORM -ldflags="-s -w" -o /app/exe ./exe/exe.go

# Compilamos el cliente
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-s -w" -o /client/client ./client/client.go

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/exe .
COPY --from=builder /app/client .

CMD ["./exe"]