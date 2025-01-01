FROM golang:1.22 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/failsafe-discord-bot cmd/discord-bot/main.go

FROM alpine

COPY --from=build /app/failsafe-discord-bot /usr/bin/

CMD [ "/usr/bin/failsafe-discord-bot" ]
