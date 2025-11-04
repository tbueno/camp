package utils

import (
	"testing"
)

func TestCommandReturn(t *testing.T) {
	tests := []struct {
		name    string
		command string
		args    []string
		want    string
		wantErr bool
	}{
		{
			name:    "Echo command",
			command: "echo",
			args:    []string{"hello"},
			want:    "hello\n",
			wantErr: false,
		},
		{
			name:    "Invalid command",
			command: "invalid_command",
			args:    []string{},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CommandReturn(tt.command, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommandReturn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CommandReturn() = %v, want %v", got, tt.want)
			}
		})
	}
}
