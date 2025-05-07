package cli

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_Run_E2E(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e testing in short mode")
	}

	type args struct {
		args []string
	}
	tests := []struct {
		name              string
		args              args
		wantExitCode      int
		wantErrorContains string
		wantJSONLRowCount int
	}{
		{
			name: "simplest json",
			args: args{
				args: []string{"-config", "../../examples/01-simplest.json"},
			},
			wantExitCode:      0,
			wantJSONLRowCount: 10,
		},
		{
			name: "simplest yaml",
			args: args{
				args: []string{"-config", "../../examples/01-simplest.yml"},
			},
			wantExitCode:      0,
			wantJSONLRowCount: 10,
		},
		{
			name: "empty config path",
			args: args{
				args: []string{"-config", ""},
			},
			wantExitCode:      1,
			wantErrorContains: "config file path is required",
		},
		{
			name: "invalid config path",
			args: args{
				args: []string{"-config", "/tmp/file-not-exists.json"},
			},
			wantExitCode:      1,
			wantErrorContains: "no such file or directory",
		},
		{
			name: "invalid param",
			args: args{
				args: []string{"-invalid-param"},
			},
			wantExitCode:      1,
			wantErrorContains: "flag provided but not defined",
		},
		{
			name: "bad extension",
			args: args{
				args: []string{"-config", "../../testdata/simplest.bad-extension"},
			},
			wantExitCode:      1,
			wantErrorContains: "unsupported config format",
		},
		{
			name: "decode issue w. enums",
			args: args{
				args: []string{"-config", "../../testdata/bad_custom_types.yml"},
			},
			wantExitCode:      1,
			wantErrorContains: "failed to decode config",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
			exitCode := Run(tt.args.args, stdout, stderr)
			assert.Equal(t, tt.wantExitCode, exitCode, "unexpected exit code")

			stderrString := stderr.String()
			assert.Contains(t, stderrString, tt.wantErrorContains)

			// Only if expecting success or line count is meaningless
			if exitCode == 0 {
				stdoutString := strings.TrimSpace(stdout.String())
				rowCount := len(strings.Split(stdoutString, "\n"))
				assert.Equal(t, tt.wantJSONLRowCount, rowCount, "row count")
			}
		})
	}
}
