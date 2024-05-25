FROM golang:1.22.3-alpine 


WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN GOOS=linux go build -o ./build/ ./cmd/orders/main.go

CMD ["./build/main"]
