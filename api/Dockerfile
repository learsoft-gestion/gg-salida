FROM golang:1.21 as sindicatos_api_builder

WORKDIR /app

# Copy go.mod and go.sum to download dependencies
# Ensure these files by running 'go mod init' and 'go mod tidy' in your project directory
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the working directory
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /sindicatos-api

# Use a Docker multi-stage build to create a lean production image.
# Start with a smaller image without the Go toolchain.
FROM alpine:latest  

# Set the working directory inside the container
WORKDIR /

# Import the compiled executable from the first stage.
COPY --from=sindicatos_api_builder /sindicatos-api /sindicatos-api

COPY . .

# Copy timezone data
RUN apk add --no-cache tzdata

# Set the timezone environment variable
ENV TZ=America/Argentina/Buenos_Aires

# Expose port (the same as in the ENV SV_ADDR)
# EXPOSE 8088

# Command to run the executable
# CMD ["/sindicatos"]
