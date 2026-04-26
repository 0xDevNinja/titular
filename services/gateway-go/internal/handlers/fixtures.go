package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/types"
)

// fixtureStore holds the in-memory fixture data loaded once at startup.
type fixtureStore struct {
	agents []types.Agent
	trades []types.Trade
	jobs   []types.Job
}

var (
	once         sync.Once
	globalStore  *fixtureStore
	globalStoreE error
)

// fixturesDir returns the path to the fixtures directory relative to this file's
// location. This works both in normal builds and in `go test` runs.
func fixturesDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "fixtures")
}

// loadFixtures reads agents.json, trades.json, and jobs.json from the fixtures directory.
func loadFixtures(dir string) (*fixtureStore, error) {
	store := &fixtureStore{}

	agentsPath := filepath.Join(dir, "agents.json")
	agentsData, err := os.ReadFile(agentsPath)
	if err != nil {
		return nil, fmt.Errorf("read agents fixture: %w", err)
	}
	if err := json.Unmarshal(agentsData, &store.agents); err != nil {
		return nil, fmt.Errorf("parse agents fixture: %w", err)
	}

	tradesPath := filepath.Join(dir, "trades.json")
	tradesData, err := os.ReadFile(tradesPath)
	if err != nil {
		return nil, fmt.Errorf("read trades fixture: %w", err)
	}
	if err := json.Unmarshal(tradesData, &store.trades); err != nil {
		return nil, fmt.Errorf("parse trades fixture: %w", err)
	}

	jobsPath := filepath.Join(dir, "jobs.json")
	jobsData, err := os.ReadFile(jobsPath)
	if err != nil {
		return nil, fmt.Errorf("read jobs fixture: %w", err)
	}
	if err := json.Unmarshal(jobsData, &store.jobs); err != nil {
		return nil, fmt.Errorf("parse jobs fixture: %w", err)
	}

	return store, nil
}

// defaultStore returns the singleton fixture store, loading from the embedded
// fixtures directory on first call.
func defaultStore() (*fixtureStore, error) {
	once.Do(func() {
		globalStore, globalStoreE = loadFixtures(fixturesDir())
	})
	return globalStore, globalStoreE
}

// NewStoreFromDir creates a fresh (non-singleton) store from the given dir.
// Used in tests to inject a specific fixtures path.
func NewStoreFromDir(dir string) (*fixtureStore, error) {
	return loadFixtures(dir)
}
