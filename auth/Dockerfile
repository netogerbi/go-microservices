FROM golang:1.18-alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o authApp ./cmd/api

RUN chmod +x authApp

# build a tiny linux image to run the app
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/authApp .

ENTRYPOINT [ "./authApp" ]