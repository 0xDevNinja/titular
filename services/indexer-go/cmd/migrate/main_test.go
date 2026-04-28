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
			code := run(tc.args, tc.dsn, emptyEnv, stdout, stderr)
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

// TestRunRefusesDestructiveByDefault makes sure "down" and "force" cannot
// be run without MIGRATE_ALLOW_DESTRUCTIVE=1. The gate fires before db.New,
// so we never need a live database for this assertion.
func TestRunRefusesDestructiveByDefault(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
	}{
		{"down", []string{"down"}},
		{"force", []string{"force", "3"}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			stdout, stderr := captureStdio(t)
			code := run(tc.args, "postgres://x:y@localhost:5432/z", emptyEnv, stdout, stderr)
			if code != 2 {
				t.Errorf("exit code = %d, want 2", code)
			}
			data, _ := io.ReadAll(reopen(t, stderr))
			if !strings.Contains(string(data), "destructive command") {
				t.Errorf("stderr = %q, want to mention destructive command", data)
			}
			if !strings.Contains(string(data), "MIGRATE_ALLOW_DESTRUCTIVE=1") {
				t.Errorf("stderr = %q, want to name the override env var", data)
			}
		})
	}
}

// TestRunAcceptsDestructiveWithEnv confirms the gate stops being hit once
// MIGRATE_ALLOW_DESTRUCTIVE=1 is set. We point at a DSN that no Postgres
// is listening on, so the call fails *past* the guard with a connection
// error (exit 1 or 2 from the migrate library) — never the guard's exit 2
// + "destructive command refused" message. We assert on the absence of
// the guard message rather than the exact downstream code.
func TestRunAcceptsDestructiveWithEnv(t *testing.T) {
	t.Parallel()

	getenv := func(k string) string {
		if k == "MIGRATE_ALLOW_DESTRUCTIVE" {
			return "1"
		}
		return ""
	}

	for _, args := range [][]string{{"down"}, {"force", "1"}} {
		args := args
		t.Run(args[0], func(t *testing.T) {
			t.Parallel()
			stdout, stderr := captureStdio(t)
			// 127.0.0.1:1 is reserved → connect fails immediately.
			_ = run(args, "postgres://x:y@127.0.0.1:1/z?sslmode=disable&connect_timeout=1", getenv, stdout, stderr)
			data, _ := io.ReadAll(reopen(t, stderr))
			if strings.Contains(string(data), "destructive command") {
				t.Errorf("guard fired with env set; stderr = %q", data)
			}
		})
	}
}

// emptyEnv simulates a process started with no environment variables. We
// use it instead of os.Getenv to keep the test hermetic — a developer
// running with MIGRATE_ALLOW_DESTRUCTIVE=1 in their shell would otherwise
// silently flip these tests' expectations.
func emptyEnv(string) string { return "" }

// TestRunHelp confirms the help command exits 0 and prints usage.
func TestRunHelp(t *testing.T) {
	t.Parallel()

	for _, arg := range []string{"-h", "--help", "help"} {
		arg := arg
		t.Run(arg, func(t *testing.T) {
			t.Parallel()
			stdout, stderr := captureStdio(t)
			code := run([]string{arg}, "postgres://x:y@localhost:5432/z", emptyEnv, stdout, stderr)
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
