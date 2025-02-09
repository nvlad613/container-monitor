FROM golang:1.23.6-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/pinger ./cmd/app/main.go

FROM alpine:latest
RUN apk --no-cache add docker-cli
WORKDIR /app

# Копируем бинарник и конфиг
COPY --from=build /app/pinger /app/pinger
COPY --from=build /app/config/application.yml /app/

RUN mkdir -p /app/logs

CMD ["/app/pinger"]