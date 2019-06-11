package service

import (
	"context"
	"reflect"
	"testing"

	v1 "github.com/maelvls/users-grpc/schema/health/v1"
)

func TestHealthImpl_Check(t *testing.T) {
	type args struct {
		ctx  context.Context
		args *v1.HealthCheckRequest
	}
	tests := []struct {
		name    string
		h       *HealthImpl
		args    args
		want    *v1.HealthCheckResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.Check(tt.args.ctx, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("HealthImpl.Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HealthImpl.Check() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHealthImpl_Watch(t *testing.T) {
	type args struct {
		in0 *v1.HealthCheckRequest
		in1 v1.Health_WatchServer
	}
	tests := []struct {
		name    string
		h       *HealthImpl
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.h.Watch(tt.args.in0, tt.args.in1); (err != nil) != tt.wantErr {
				t.Errorf("HealthImpl.Watch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
