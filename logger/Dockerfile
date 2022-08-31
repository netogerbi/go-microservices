FROM golang:1.18-alpine as builder

WORKDIR /app

COPY . .

# once this is production container we can compile protobuf on build
# and then comile golang code
RUN apt update
RUN apt install -y protobuf-compiler
RUN GO111MODULE=on
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

RUN CGO_ENABLED=0 go build -o loggerApp ./cmd/api

RUN chmod +x loggerApp

# build a tiny linux image to run the app
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/loggerApp .

ENTRYPOINT [ "./loggerApp" ]