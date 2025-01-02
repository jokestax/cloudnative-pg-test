FROM golang:alpine3.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o main .

FROM alpine:3.21

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]