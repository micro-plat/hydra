package pub

import (
	"testing"

	"golang.org/x/crypto/ssh"
)

func Test_sshClient_Bind(t *testing.T) {
	type fields struct {
		ip          string
		userName    string
		pwd         string
		client      *ssh.Client
		tmpDir      string
		tmpPath     string
		tmpFile     string
		localPath   string
		projectPath string
	}
	type args struct {
		host      string
		localpath string
		pwd       string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sshClient{
				ip:          tt.fields.ip,
				userName:    tt.fields.userName,
				pwd:         tt.fields.pwd,
				client:      tt.fields.client,
				tmpDir:      tt.fields.tmpDir,
				tmpPath:     tt.fields.tmpPath,
				tmpFile:     tt.fields.tmpFile,
				localPath:   tt.fields.localPath,
				projectPath: tt.fields.projectPath,
			}
			if err := s.Bind(tt.args.host, tt.args.localpath, tt.args.pwd); (err != nil) != tt.wantErr {
				t.Errorf("sshClient.Bind() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
