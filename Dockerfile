# ---------- BUILD STAGE ----------
FROM golang:1.24-alpine AS build
WORKDIR /app

# Copy go mod files and download dependencies first (cached layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the project
COPY . .

# Build the Go binary
RUN go build -o server ./cmd

# ---------- RUNTIME STAGE ----------
FROM alpine:latest
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/server .

# Expose the service port (same as in your code)
EXPOSE 8080

# Run the app
CMD ["./server"]
