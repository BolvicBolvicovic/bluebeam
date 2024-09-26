FROM golang:1.23

# Set the working directory
WORKDIR /app

# Set the Go proxy to direct
RUN go env -w GOPROXY=direct

# Copy go.mod and go.sum (if you have it) first to leverage Docker cache
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o app ./cmd/main.go

EXPOSE 3333

# Command to run the application
CMD ["./app"]
