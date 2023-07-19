# Authentication Service

This project represents a service for authenticating users.

## Installation

### Running with Go

To run the service locally with Go, execute the following commands:

```bash
go get
go run main.go
```

### Running with Docker

To run the service with Docker, you first need to build a Docker image. Execute the following commands:

```bash
docker build -t authorization:latest .
```

Then you can run the service with the following command:

```bash
docker run -p 8020:8020 authorization:latest
```

## Uploading Docker Image to Docker Hub

If you want to upload your Docker image to Docker Hub, execute the following commands:

```bash
docker login
docker tag authorization:latest uruk/authorization:latest
docker push uruk/authorization:latest
```

## Usage

The service listens on port 8020 and provides an API for user authentication.

## Configuration

The service uses a `config.yaml` file for configuration. In this file, you can specify parameters for database connection, JWT secret, and other settings.



git remote:
```
git remote set-url origin https://UrukHan:<your_token>@github.com/UrukHan/Auth-GO.git
```

