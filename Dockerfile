# STAGE 1: The "builder" stage
# Use the official Go image. Using alpine for a smaller build stage.
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of your source code
COPY . .

# Build a static, CGO-disabled binary for Linux
# This is crucial for a minimal "scratch" or "alpine" final image
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/main .

# ---

# STAGE 2: The "final" stage
# Start from a minimal-footprint image
FROM alpine:latest
WORKDIR /app

# Copy *only* the compiled binary from the "builder" stage
COPY --from=builder /app/main .

# Set a default PORT. Cloud Run will override this at runtime.
# This is just a default for local testing.
ENV PORT 8080
EXPOSE 8080

# This is the command that will run your application
# ./main is the binary we built in STAGE 1
CMD ["./main"]