# Use the official Golang image as the base image
FROM golang:1.22.4

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application code to the working directory
COPY . .

# Download and install any required Go modules
RUN go mod download

# Build the Go application
RUN go build -o main .

# Expose the port that your Go application runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
