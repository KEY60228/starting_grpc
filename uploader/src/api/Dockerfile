FROM golang:1.17-alpine3.14

RUN apk add --update --no-cache curl make protobuf protobuf-dev git grpc-cli
RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go get -u github.com/golang/protobuf/protoc-gen-go go.uber.org/zap github.com/grpc-ecosystem/go-grpc-middleware

WORKDIR /golang