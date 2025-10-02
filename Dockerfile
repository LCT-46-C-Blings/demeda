FROM golang:latest AS builder

RUN go install github.com/swaggo/swag/cmd/swag@latest

# Fetch dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN make

# Create final image
FROM alpine
WORKDIR /
COPY --from=builder /build/demeda .
EXPOSE 8080
CMD ["./demeda"]