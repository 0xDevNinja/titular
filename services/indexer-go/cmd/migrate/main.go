// Command migrate applies, reverts, and inspects the indexer Postgres schema.
//
// Usage:
//
//	migrate up                # apply all pending migrations
//	migrate down              # revert every applied migration
//	migrate steps <n>         # apply n (n>0) or revert -n (n<0) migrations
//	migrate version           # print current version + dirty flag
//	migrate force <v>         # set version to v without running SQL (recovery)
//
// The DSN is read from $DATABASE_URL. Both the libpq form
// (postgres://user:pw@host:port/db) and the migrate-specific pgx5:// form
// are accepted.
//
// Exit codes:
//
//	0  success (including "no change" results)
//	1  bad arguments / config
//	2  migration failure (the database may be left in a dirty state — run
//	   `migrate version` and consult the runbook before forcing)
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/0xDevNinja/titular/services/indexer-go/internal/db"
)

const usage = `migrate — apply or revert indexer-go Postgres migrations

Usage:
  migrate <command> [args]

Commands:
  up              apply every pending migration
  down            revert every applied migration
  steps <n>       apply n (n>0) or revert -n (n<0) migrations
  version         print current schema version and dirty flag
  force <v>       set schema_migrations to (v, dirty=false) without running SQL

Environment:
  DATABASE_URL    Postgres DSN (postgres://… or pgx5://…). Required.
`

func main() {
	// Use a custom flagset so -h prints our combined usage rather than the
	// default Go flag dump.
	fs := flag.NewFlagSet("migrate", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() { fmt.Fprint(os.Stderr, usage) }

	if err := fs.Parse(os.Args[1:]); err != nil {
		os.Exit(1)
	}
	args := fs.Args()
	if len(args) == 0 {
		fs.Usage()
		os.Exit(1)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		fmt.Fprintln(os.Stderr, "migrate: DATABASE_URL is not set")
		os.Exit(1)
	}

	code := run(args, dsn, os.Stdout, os.Stderr)
	os.Exit(code)
}

// action is a parsed CLI invocation. We validate args before opening any
// network connection so misuse fails fast with a clear message instead of
// hiding behind a connect-refused error.
type action struct {
	cmd string // "up", "down", "steps", "version", "force"
	n   int    // arg for steps / force
}

// parseArgs validates the command + its operand. The ok flag distinguishes
// "this is a help request" (exit 0) from "this is a usage error" (exit 1)
// without making the caller stuff that into the action struct.
func parseArgs(args []string, stdout, stderr *os.File) (a action, code int, done bool) {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return action{}, 1, true
	}
	cmd := args[0]
	switch cmd {
	case "up", "down", "version":
		return action{cmd: cmd}, 0, false
	case "steps":
		if len(args) < 2 {
			fmt.Fprintln(stderr, "migrate: steps: missing <n>")
			return action{}, 1, true
		}
		n, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Fprintf(stderr, "migrate: steps: %v\n", err)
			return action{}, 1, true
		}
		return action{cmd: cmd, n: n}, 0, false
	case "force":
		if len(args) < 2 {
			fmt.Fprintln(stderr, "migrate: force: missing <v>")
			return action{}, 1, true
		}
		v, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Fprintf(stderr, "migrate: force: %v\n", err)
			return action{}, 1, true
		}
		return action{cmd: cmd, n: v}, 0, false
	case "-h", "--help", "help":
		fmt.Fprint(stdout, usage)
		return action{}, 0, true
	default:
		fmt.Fprintf(stderr, "migrate: unknown command %q\n\n", cmd)
		fmt.Fprint(stderr, usage)
		return action{}, 1, true
	}
}

// run is split out so tests can drive the dispatcher without execing a
// subprocess. It returns the exit code instead of calling os.Exit.
func run(args []string, dsn string, stdout, stderr *os.File) int {
	a, code, done := parseArgs(args, stdout, stderr)
	if done {
		return code
	}

	mig, err := db.New(dsn)
	if err != nil {
		fmt.Fprintf(stderr, "migrate: open: %v\n", err)
		return 1
	}
	defer func() {
		if cerr := mig.Close(); cerr != nil {
			fmt.Fprintf(stderr, "migrate: close: %v\n", cerr)
		}
	}()

	switch a.cmd {
	case "up":
		if err := mig.Up(); err != nil {
			fmt.Fprintf(stderr, "migrate: up: %v\n", err)
			return 2
		}
		fmt.Fprintln(stdout, "ok: up")
		return 0

	case "down":
		if err := mig.Down(); err != nil {
			fmt.Fprintf(stderr, "migrate: down: %v\n", err)
			return 2
		}
		fmt.Fprintln(stdout, "ok: down")
		return 0

	case "steps":
		if err := mig.Steps(a.n); err != nil {
			fmt.Fprintf(stderr, "migrate: steps: %v\n", err)
			return 2
		}
		fmt.Fprintf(stdout, "ok: steps %d\n", a.n)
		return 0

	case "version":
		v, dirty, err := mig.Version()
		if err != nil {
			fmt.Fprintf(stderr, "migrate: version: %v\n", err)
			return 2
		}
		fmt.Fprintf(stdout, "version=%d dirty=%t\n", v, dirty)
		return 0

	case "force":
		if err := mig.Force(a.n); err != nil {
			fmt.Fprintf(stderr, "migrate: force: %v\n", err)
			return 2
		}
		fmt.Fprintf(stdout, "ok: force %d\n", a.n)
		return 0
	}

	// Unreachable: parseArgs only returns the cases above.
	fmt.Fprintf(stderr, "migrate: internal error: unhandled command %q\n", a.cmd)
	return 1
}
