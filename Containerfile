# Use the Red Hat UBI 9 Go Toolset base image
FROM registry.access.redhat.com/ubi9/go-toolset:latest as builder

USER root

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY app .

# Build the Go application
RUN go mod init app && \
    go mod tidy && \
    go build -o main .

# Start a new stage for the minimal runtime image
FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built executable from the builder stage
COPY --from=builder /app/main .

# Set the CMD to run the application with environment variables for FILE, KEY, and RUN_TYPE
CMD ["./main"]
