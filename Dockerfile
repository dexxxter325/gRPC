FROM golang:1.22.0-alpine3.19 AS builder

RUN go version

COPY ./ /GRPC

WORKDIR /GRPC

RUN go mod download
RUN go build -o GRPC ./cmd/app/main.go

EXPOSE 8000

EXPOSE 8080

CMD ["./GRPC"]