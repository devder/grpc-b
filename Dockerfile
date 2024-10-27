# use multi-stage builds to reduce the size of the final image

# build stage
FROM golang:1.23.2-alpine3.20 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o main main.go
RUN apk add --no-cache curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz


# run stage
FROM alpine:3.20
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/main .
COPY --from=builder /usr/src/app/migrate .
COPY db/migrations ./db/migrations

# TODO: remove
COPY app.env .
COPY start.sh .

EXPOSE 8080
CMD ["/usr/src/app/main"]
ENTRYPOINT ["/usr/src/app/start.sh"]
