package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kamilch1k/shardlab/internal/sim"
)

func TestSimulateEndpoint(t *testing.T) {
	handler := NewHandler()
	body, err := json.Marshal(scenario())
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/simulate", bytes.NewReader(body))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with %s", recorder.Code, recorder.Body.String())
	}
	var report sim.Report
	if err := json.NewDecoder(recorder.Body).Decode(&report); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if report.Keys == 0 || report.MovementRatio == 0 {
		t.Fatalf("unexpected report: %#v", report)
	}
}

func TestStreamEndpointEmitsReportEvent(t *testing.T) {
	handler := NewHandler()
	body, _ := json.Marshal(scenario())
	request := httptest.NewRequest(http.MethodPost, "/api/simulate/stream", bytes.NewReader(body))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), "event: report") {
		t.Fatalf("missing report event: %s", recorder.Body.String())
	}
}

func scenario() sim.Scenario {
	return sim.Scenario{
		Name:        "api-test",
		Algorithm:   sim.AlgorithmConsistent,
		NodesBefore: []sim.Node{{ID: "a"}, {ID: "b"}, {ID: "c"}},
		NodesAfter:  []sim.Node{{ID: "a"}, {ID: "b"}, {ID: "c"}, {ID: "d"}},
		Workload: sim.Workload{
			KeyCount: 1000,
			Requests: 4000,
			Seed:     11,
		},
		Ring: sim.RingConfig{VirtualNodes: 128},
	}
}
