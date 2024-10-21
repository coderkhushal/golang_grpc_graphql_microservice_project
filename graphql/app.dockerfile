FROM golang:1.23.1-alpine AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/coderkhushal/go-grpc-graphql-microservices
COPY go.mod go.sum ./
COPY vendor vendor
COPY account account
COPY catalog catalog
COPY order order
COPY graphql graphql
RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./graphql

FROM alpine
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]