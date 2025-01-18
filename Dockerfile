# use multi-stage builds to reduce the size of the final image

# build stage
FROM golang:1.23-alpine3.21 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o main main.go
# RUN apk add --no-cache curl
# RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz


# run stage
FROM alpine:3.21
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/main .
# COPY --from=builder /usr/src/app/migrate .
COPY db/migrations ./db/migrations

# app.env for prod is injected in CI workflow
COPY app.env .
COPY start.sh .

EXPOSE 8080
CMD ["/usr/src/app/main"]
ENTRYPOINT ["/usr/src/app/start.sh"]
