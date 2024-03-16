FROM golang as build

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o failsafe cmd/discord-bot/main.go

ENTRYPOINT [ "./failsafe" ]
