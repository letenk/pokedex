# Build stage
FROM golang:1.18-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY ./scripts/start.sh .
COPY ./scripts/wait-for.sh .

EXPOSE 3000
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]