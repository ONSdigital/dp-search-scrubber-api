# Use a base image with the latest version of Go
FROM golang:1.23.4-bullseye

# Set the working directory to /app
WORKDIR /app

# Copy the code into the container
COPY . .

# Build the Go binary
RUN go build -o main .

# Expose port 28700
EXPOSE 28700

# Run the Go binary when the container starts
CMD ["./main"]