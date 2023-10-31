# Use a base image with Go installed
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the go source files and other necessary files
COPY . .

# Build your tool
RUN go build -o break-check

# Set the entrypoint to your tool
ENTRYPOINT ["/app/break-check"]
