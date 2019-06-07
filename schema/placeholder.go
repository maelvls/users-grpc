//go:generate go get -u github.com/golang/protobuf/protoc-gen-go
//go:generate protoc quote.proto --go_out=plugins=grpc:./quote

// The purpose of this file is only to hold the //go:generate lines.

package schema
