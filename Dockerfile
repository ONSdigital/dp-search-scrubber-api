# Use a base image with the latest version of Go
FROM golang:latest

# Set the working directory to /app
WORKDIR /app

# Copy the code into the container
COPY . .

# Build the Go binary
RUN go build -o main .

# Expose port 5000
EXPOSE 5000

# Run the Go binary when the container starts
CMD ["./main"]