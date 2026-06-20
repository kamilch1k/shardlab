package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIProducesTextReport(t *testing.T) {
	dir := t.TempDir()
	scenario := filepath.Join(dir, "scenario.json")
	if err := os.WriteFile(scenario, []byte(`{
  "name": "test",
  "algorithm": "consistent",
  "nodesBefore": [{"id":"a"},{"id":"b"},{"id":"c"}],
  "nodesAfter": [{"id":"a"},{"id":"b"},{"id":"c"},{"id":"d"}],
  "workload": {"keyCount": 1000, "requests": 4000, "seed": 7},
  "ring": {"virtualNodes": 128}
}`), 0o600); err != nil {
		t.Fatalf("write scenario: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := run([]string{"-input", scenario}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "movement=") || !strings.Contains(stdout.String(), "after_loads") {
		t.Fatalf("unexpected output: %s", stdout.String())
	}
}
