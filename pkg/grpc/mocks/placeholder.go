//go:generate go run -mod=mod github.com/golang/mock/mockgen -build_flags=-mod=mod -package mocks -destination ./mock_service.go -source=../user.go

// The purpose of this file is only to hold the //go:generate lines.

package mocks
