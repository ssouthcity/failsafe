FROM golang:1.22 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /app/discord-bot cmd/discord-bot/main.go

FROM scratch

COPY --from=build /app/discord-bot /usr/bin/

CMD [ "discord-bot" ]
