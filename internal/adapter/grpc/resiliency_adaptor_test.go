package grpc

import (
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGenerateErrStatus(t *testing.T) {
	tests := []struct {
		name       string
		statusCode uint32
		wantCode   codes.Code
		wantErr    bool
	}{
		{name: "OK", statusCode: 0, wantCode: codes.OK},
		{name: "Unavailable", statusCode: 14, wantCode: codes.Unavailable, wantErr: true},
		{name: "unknown numeric status", statusCode: 99, wantCode: codes.Internal, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := generateErrStatus(tt.statusCode)
			if !tt.wantErr {
				if err != nil {
					t.Fatalf("generateErrStatus(%d) = %v, want nil", tt.statusCode, err)
				}
				return
			}

			if err == nil {
				t.Fatalf("generateErrStatus(%d) = nil, want error", tt.statusCode)
			}
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("status.Code(err) = %v, want %v", got, tt.wantCode)
			}
		})
	}
}
