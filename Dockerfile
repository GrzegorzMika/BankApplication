FROM golang:1.20-alpine3.18 AS builder

WORKDIR /app

COPY . .

RUN go install github.com/gobuffalo/pop/v6/soda@latest

RUN go build -o main main.go

FROM alpine:3.18 AS app

WORKDIR /app

COPY ./db/migrations ./db/database.yml /db/
COPY ./start.sh .
COPY --from=builder /app/main .
COPY --from=builder /app/app.env .
COPY --from=builder /go/bin/soda ./soda

RUN chmod a+x /app/soda

EXPOSE 8080

CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]