package mqc

import "testing"

func TestServer_Start(t *testing.T) {
	tests := []struct {
		name    string
		s       *Server
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Start(); (err != nil) != tt.wantErr {
				t.Errorf("Server.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
