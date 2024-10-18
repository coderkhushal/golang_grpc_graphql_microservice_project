# Build Stage
FROM golang:1.13-alpine3.11 AS build
# Install necessary build tools
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/akhilsharma90/go-graphql-microservice

# Copy Go mod and sum files first to cache dependencies (faster builds)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY vendor vendor
COPY order order

# Build the Go application
RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./order/cmd/order

# Runtime Stage
FROM alpine:3.11
WORKDIR /usr/bin

# Copy the built application from the build stage
COPY --from=build /go/bin/app .

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./app"]
