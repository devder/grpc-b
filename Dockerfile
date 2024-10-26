# use multi-stage builds to reduce the size of the final image

# build stage
FROM golang:1.23.2-alpine3.20 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o main main.go

# run stage
FROM alpine:3.20
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/main .

# TODO: remove
COPY app.env .

EXPOSE 8080
CMD ["/usr/src/app/main"]
