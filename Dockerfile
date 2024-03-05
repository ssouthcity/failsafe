FROM golang as build

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o failsafe cmd/dbot/main.go

ENTRYPOINT [ "./failsafe" ]
