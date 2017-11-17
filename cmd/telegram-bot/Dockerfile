FROM golang:alpine as builder
RUN apk add --no-cache git
ADD . /go/src/github.com/rpeshkov/slovnik
WORKDIR /go/src/github.com/rpeshkov/slovnik
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN go build -o app github.com/rpeshkov/slovnik/cmd/telegram-bot

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /go/src/github.com/rpeshkov/slovnik/app /app
ADD ./cmd/telegram-bot/templates /templates
ENTRYPOINT /app
EXPOSE 8080