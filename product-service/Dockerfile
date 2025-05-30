# --- Builder Stage ---
FROM golang:1.24.3-alpine AS builder

# Set the working directory inside the container
WORKDIR /product-service

COPY go.mod .
COPY go.sum .

# Download Go modules. This step is cached.
RUN go mod download

# Copy the rest of the application
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# --- Runner Stage ---
# Use a minimal Alpine Linux image for the final production image
FROM alpine:latest

# Create a non-root user and group for security best practices
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Set the working directory for the application
WORKDIR /product-service

# Copy the compiled binary from the 'builder' stage
COPY --from=builder /product-service/main .

# Just for documentation
EXPOSE 8081

# Command to run the application
# The application will listen on the port defined by the PORT environment variable.
CMD ["./main"]