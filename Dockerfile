FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go build -o ip-server main.go

EXPOSE 8081

CMD ["./ip-server"]
