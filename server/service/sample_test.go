// Putting some sample data.

package service

import "testing"

func TestUserImpl_LoadSampleUsers(t *testing.T) {
	tests := []struct {
		name    string
		svc     *UserImpl
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.svc.LoadSampleUsers(); (err != nil) != tt.wantErr {
				t.Errorf("UserImpl.LoadSampleUsers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
