
#build stage
FROM golang:alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN apk add --no-cache git
RUN go get -d -v ./...
RUN go build -v -o ./build/ ./cmd/bankingd

#final stage
FROM alpine:latest


RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/app/build/bankingd /bankingd
ENTRYPOINT ./bankingd

LABEL Name=stone-desafio-go Version=0.0.1

EXPOSE 80
