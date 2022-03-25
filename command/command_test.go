package command

import (
	"testing"
)

func TestIsBinaryFile(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "/bin/ls",
			args: args{
				file: "/bin/ls",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "/etc/passwd",
			args: args{
				file: "/etc/passwd",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "/etc/not_exist_file",
			args: args{
				file: "/etc/not_exist_file",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsBinaryFile(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsBinaryFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsBinaryFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasExecPermition(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "/bin/ls",
			args: args{
				path: "/bin/ls",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "/etc/passwd",
			args: args{
				path: "/etc/passwd",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "/etc/not_exist_file",
			args: args{
				path: "/etc/not_exist_file",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HasExecPermition(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("HasExecPermition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HasExecPermition() = %v, want %v", got, tt.want)
			}
		})
	}
}
