# STAGE 1: The "builder" stage
# Use the Go version from your go.mod
FROM golang:1.24-alpine AS builder

# Set a clean working directory
WORKDIR /app

# Copy module files first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all your source code into the /app directory
# (e.g., /app/cmd, /app/internal, etc.)
COPY . .

# --- THIS IS THE KEY FIX ---
# Build the package located at ./cmd (which contains main.go)
# The output binary will be at /app/main
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/main ./cmd

# ---

# STAGE 2: The "final" stage
# Use a minimal, non-root, secure image
FROM gcr.io/distroless/static-debian11
WORKDIR /

# Copy *only* the compiled binary from the builder stage
COPY --from=builder /app/main .

# Set the default port Cloud Run will use.
# Your code MUST read this env variable.
ENV PORT="8080"
EXPOSE 8080

# Run the binary
CMD ["/main"]