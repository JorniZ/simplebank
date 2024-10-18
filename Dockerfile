# build stage
FROM golang:1.23-alpine3.20 AS builder 
WORKDIR /simple-bank-api
COPY . .
RUN go build -o main main.go

# run stage
FROM alpine:3.20
WORKDIR /simple-bank-api
COPY --from=builder /simple-bank-api/main . 
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration

EXPOSE 8080
CMD [ "/simple-bank-api/main" ]
ENTRYPOINT [ "/simple-bank-api/start.sh" ]