package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// TestRunRejectsBadInput exercises the CLI dispatcher's argument-handling
// paths without ever touching a database. We rely on the fact that db.New
// rejects an empty DSN — but the CLI rejects empty / unknown commands and
// missing operands earlier, which is what we are checking here.
func TestRunRejectsBadInput(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		args     []string
		dsn      string
		wantCode int
		wantMsg  string // substring expected in stderr
	}{
		{
			name:     "unknown command",
			args:     []string{"sideways"},
			dsn:      "postgres://x:y@localhost:5432/z",
			wantCode: 1,
			wantMsg:  `unknown command "sideways"`,
		},
		{
			name:     "steps without n",
			args:     []string{"steps"},
			dsn:      "postgres://x:y@localhost:5432/z",
			wantCode: 1,
			wantMsg:  `missing <n>`,
		},
		{
			name:     "steps with non-numeric n",
			args:     []string{"steps", "twelve"},
			dsn:      "postgres://x:y@localhost:5432/z",
			wantCode: 1,
			wantMsg:  `steps`,
		},
		{
			name:     "force without v",
			args:     []string{"force"},
			dsn:      "postgres://x:y@localhost:5432/z",
			wantCode: 1,
			wantMsg:  `missing <v>`,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			stdout, stderr := captureStdio(t)
			code := run(tc.args, tc.dsn, stdout, stderr)
			if code != tc.wantCode {
				t.Errorf("exit code = %d, want %d", code, tc.wantCode)
			}
			data, _ := io.ReadAll(reopen(t, stderr))
			if !strings.Contains(string(data), tc.wantMsg) {
				t.Errorf("stderr = %q, want to contain %q", data, tc.wantMsg)
			}
		})
	}
}

// TestRunHelp confirms the help command exits 0 and prints usage.
func TestRunHelp(t *testing.T) {
	t.Parallel()

	for _, arg := range []string{"-h", "--help", "help"} {
		arg := arg
		t.Run(arg, func(t *testing.T) {
			t.Parallel()
			stdout, stderr := captureStdio(t)
			code := run([]string{arg}, "postgres://x:y@localhost:5432/z", stdout, stderr)
			if code != 0 {
				t.Errorf("exit code = %d, want 0", code)
			}
			data, _ := io.ReadAll(reopen(t, stdout))
			if !strings.Contains(string(data), "Commands:") {
				t.Errorf("stdout missing usage; got %q", data)
			}
		})
	}
}

// captureStdio returns a pair of *os.File temp files acting as stdout and
// stderr replacements. The dispatcher writes to *os.File so we cannot just
// pass a bytes.Buffer.
func captureStdio(t *testing.T) (stdout, stderr *os.File) {
	t.Helper()

	out, err := os.CreateTemp(t.TempDir(), "stdout-*")
	if err != nil {
		t.Fatalf("create stdout: %v", err)
	}
	t.Cleanup(func() { _ = out.Close() })

	er, err := os.CreateTemp(t.TempDir(), "stderr-*")
	if err != nil {
		t.Fatalf("create stderr: %v", err)
	}
	t.Cleanup(func() { _ = er.Close() })

	return out, er
}

// reopen rewinds the file to start so the caller can read what was written.
func reopen(t *testing.T, f *os.File) io.Reader {
	t.Helper()
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatalf("seek: %v", err)
	}
	// Read everything into memory — these tests only emit a few bytes.
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, f); err != nil {
		t.Fatalf("read: %v", err)
	}
	return &buf
}
