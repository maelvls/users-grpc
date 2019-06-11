//go:generate go get -u github.com/golang/protobuf/protoc-gen-go
//go:generate protoc user.proto --go_out=plugins=grpc:./user
//go:generate protoc health.proto --go_out=plugins=grpc:./health/v1

// The purpose of this file is only to hold the //go:generate lines.

package schema
