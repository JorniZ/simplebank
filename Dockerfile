# build stage
FROM golang:1.23-alpine3.20 AS builder 
WORKDIR /simple-bank-api
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz

# run stage
FROM alpine:3.20
WORKDIR /simple-bank-api
COPY --from=builder /simple-bank-api/main .
COPY --from=builder /simple-bank-api/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration

EXPOSE 8080
CMD [ "/simple-bank-api/main" ]
ENTRYPOINT [ "/simple-bank-api/start.sh" ]