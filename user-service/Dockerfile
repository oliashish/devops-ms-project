# --- Builder Stage ---
FROM golang:1.24.3-alpine AS builder

# Set the working directory inside the container
WORKDIR /user-service

COPY go.mod .
COPY go.sum .

# Download Go modules. This step is cached.
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the Go application
# CGO_ENABLED=0: Disables CGo, which is important for creating statically linked binaries.
# GOOS=linux: Specifies the target operating system.
# -a -installsuffix cgo: Forces rebuilding of packages that use CGo and adds a suffix to avoid conflicts.
# -o main: Specifies the output binary name as 'main'.
# .: Compiles the current directory.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# --- Runner Stage ---
# Use a minimal Alpine Linux image for the final production image
FROM alpine:latest

# Create a non-root user and group for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Set the working directory for the application
WORKDIR /user-service

# Copy the compiled binary from the 'builder' stage
COPY --from=builder /user-service/main .

# just for Documentation
EXPOSE 8080

# Command to run the application
# The application will listen on the port defined by the PORT environment variable.
CMD ["./main"]