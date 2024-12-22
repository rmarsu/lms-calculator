FROM golang:1.23-alpine
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o myapp cmd/main.go

CMD ["./myapp"]
