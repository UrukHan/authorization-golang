# Начинаем с официального образа Go
FROM golang:1.16-alpine AS build

# Установка зависимостей для gcc
RUN apk --no-cache add gcc g++ make

# Установка рабочей директории в контейнере
WORKDIR /app

# Копируем go mod и sum файлы
COPY go.mod go.sum ./

# Скачиваем все зависимости. Зависимости будут кэшироваться, если go.mod и go.sum не изменяются
RUN go mod download

# Копируем исходный код в контейнер
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Начинаем новую стадию сборки для создания минимального образа
FROM alpine:3.14

# Настраиваем рабочую директорию
WORKDIR /root/

# Копирование исполняемого файла из предыдущей стадии
COPY --from=build /app/main .

# Экспонируем порт, на котором ваше приложение будет работать
EXPOSE 8020

# Запуск приложения
CMD ["./main"]
