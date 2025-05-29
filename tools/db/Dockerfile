# Use the official Go image as the base image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Install necessary dependencies
RUN apk add --no-cache ca-certificates

# Copy the Go script and any required files into the container
COPY db_cleanup.go ./

# Copy SSL certificates if needed
# COPY /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Build the Go script
RUN go build -o db_cleanup db_cleanup.go

# Set environment variables (can be overridden at runtime)
ENV BASE_URL=https://localhost:8081
ENV HOURS_TO_SUBTRACT=120
ENV TIMEZONE=UTC
ENV SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt

# Command to run the script
CMD ["./db_cleanup"]