//go:generate go get -u github.com/golang/protobuf/protoc-gen-go
//go:generate protoc quote.proto --go_out=plugins=grpc:./quote
//go:generate protoc health.proto --go_out=plugins=grpc:./health/v1

// The purpose of this file is only to hold the //go:generate lines.

package schema
