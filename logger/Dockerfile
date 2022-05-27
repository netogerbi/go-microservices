FROM golang:1.18-alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

RUN chmod +x brokerApp

# build a tiny linux image to run the app
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/brokerApp .

ENTRYPOINT [ "./brokerApp" ]