# Start with the Go official image
FROM golang:1.21-rc-alpine3.17 AS build

# Set the working directory in the container
WORKDIR /app

# Install ca-certificates
RUN apk --no-cache add ca-certificates

# Install git
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Set environment variable to bypass the proxy for github.com
ENV GOPRIVATE=github.com

# Download all dependencies. They will be cached if go.mod and go.sum files do not change
RUN go mod download
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Set environment variables
ENV ADMIN_PASSWORD=
ENV ADMIN_EMAIL=
ENV JWTSECRET=
ENV DB_HOST=
ENV DB_USER=
ENV DB_DBNAME=
ENV DB_PASSWORD=
ENV DB_SSLMODE=

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

##################
# Start a new stage from scratch to create a minimal image
FROM scratch

# Copy the binary from the previous stage
COPY --from=build /app/main .

# Copy the config file from the previous stage
COPY --from=build /app/config.yaml .

# Expose the port your app runs on
EXPOSE 8020

# Run the application
CMD ["./main"]
