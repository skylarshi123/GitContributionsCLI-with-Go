# Start from the official Golang image
FROM golang:1.23

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum first (if they exist)
# This is done separately to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
# Using go mod download instead of go get for reproducible builds
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# CGO_ENABLED=0 builds a statically linked binary
# -o specifies the output binary name
RUN CGO_ENABLED=0 go build -o gitcontrib

# Create a directory to store git repositories
RUN mkdir /repos

# Final running stage
# Command to run when container starts
# Note: this is just a default, can be overridden at runtime
ENTRYPOINT ["./gitcontrib"]

# Default arguments (can be overridden at runtime)
CMD ["-email", "default@example.com"]

# Docker commands for building and running:
# Build:
#   docker build -t gitcontrib .
#
# Run (example commands):
# 1. Basic run with email:
#   docker run -it gitcontrib -email "your@email.com"
#
# 2. Mount local git repos and specify email:
#   docker run -it -v /path/to/your/repos:/repos gitcontrib -add /repos
#   docker run -it -v /path/to/your/repos:/repos gitcontrib -email "your@email.com"
#
# 3. Mount home directory for config file:
#   docker run -it -v $HOME:/root gitcontrib -email "your@email.com"
#
# 4. Complete setup with both mounts:
#   docker run -it \
#     -v $HOME:/root \
#     -v /path/to/your/repos:/repos \
#     gitcontrib -email "your@email.com"