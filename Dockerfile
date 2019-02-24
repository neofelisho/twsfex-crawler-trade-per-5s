FROM golang:latest as builder

WORKDIR /go/src/app
COPY . .

ENV GO111MODULE=on
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -o app

FROM alpine:latest as final

#Reference: `https://wiki.alpinelinux.org/wiki/Setting_the_timezone`
RUN apk add --no-cache tzdata

WORKDIR /root
COPY --from=builder /go/src/app .

ENV dataSource="http://www.twse.com.tw/en/exchangeReport/MI_5MINS?response=csv&date="
ENV apiUrl="http://mongo-api:8080/daily"
ENTRYPOINT ["./app"]

MAINTAINER neofelisho@gmail.com
