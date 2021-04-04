FROM golang:1.16.3-alpine3.13

WORKDIR ./bot

RUN go mod tidy

EXPOSE 8000

CMD ["go run bot"]