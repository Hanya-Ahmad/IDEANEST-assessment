FROM golang:1.20.5 

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=1 go build -o server ./cmd/main.go

CMD [ "./server" ]

