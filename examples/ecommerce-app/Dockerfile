# Builder
FROM golang:1.21.11-alpine AS go

MAINTAINER Vicky Phang <vickyphang11@gmail.com>

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ecommerce-app


FROM alpine:3.14

WORKDIR /app

COPY --from=go /app/ecommerce-app /app/ecommerce-app

EXPOSE 8080

CMD ["/app/ecommerce-app"]