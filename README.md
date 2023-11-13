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

To run the service with Docker, you need to build a Docker image.
It makes use of environment variables which you can pass using the --build-arg flag as shown:

```bash
docker build \
--build-arg admin_password=${ADMIN_PASSWORD} \
--build-arg admin_email=${ADMIN_EMAIL} \
--build-arg jwtSecret=${JWTSECRET} \
--build-arg db_host=${DB_HOST} \
--build-arg db_user=${DB_USER} \
--build-arg db_dbname=${DB_DBNAME} \
--build-arg db_password=${DB_PASSWORD} \
--build-arg db_sslmode=${DB_SSLMODE} \
-t authorization:latest .

docker build --build-arg admin_password=${ADMIN_PASSWORD} --build-arg admin_email=${ADMIN_EMAIL} --build-arg jwtSecret=${JWTSECRET} --build-arg db_host=${DB_HOST} --build-arg db_user=${DB_USER} --build-arg db_dbname=${DB_DBNAME} --build-arg db_password=${DB_PASSWORD} --build-arg db_sslmode=${DB_SSLMODE} -t authorization:latest .
```

Then you can run the service with the following command:

```bash
docker run -d --env-file .env -p 8020:8020 --name authorization_container authorization:latest
```

## Uploading Docker Image to Docker Hub

If you want to upload your Docker image to Docker Hub, execute the following commands:

```bash
docker login
docker build -t uruk/authorization-golang .
docker tag uruk/authorization-golang uruk/authorization:latest
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



sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl
