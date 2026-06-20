package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/kamilch1k/shardlab/internal/sim"
)

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("POST /api/simulate", simulate)
	mux.HandleFunc("POST /api/simulate/stream", stream)
	return mux
}

func health(writer http.ResponseWriter, _ *http.Request) {
	writeJSON(writer, http.StatusOK, map[string]string{"status": "ok", "service": "shardlab"})
}

func simulate(writer http.ResponseWriter, request *http.Request) {
	scenario, ok := decodeScenario(writer, request)
	if !ok {
		return
	}
	report, err := sim.Simulate(scenario)
	if err != nil {
		writeJSON(writer, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(writer, http.StatusOK, report)
}

func stream(writer http.ResponseWriter, request *http.Request) {
	scenario, ok := decodeScenario(writer, request)
	if !ok {
		return
	}

	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")

	steps := []struct {
		Name string `json:"step"`
		Text string `json:"message"`
	}{
		{"workload", "generated deterministic request workload"},
		{"before", "assigned workload to original topology"},
		{"after", "assigned workload to changed topology"},
		{"metrics", "computed movement, skew, and Gini metrics"},
	}

	for index, step := range steps {
		writeSSE(writer, "progress", map[string]any{
			"index":   index + 1,
			"total":   len(steps) + 1,
			"step":    step.Name,
			"message": step.Text,
		})
		time.Sleep(20 * time.Millisecond)
	}

	report, err := sim.Simulate(scenario)
	if err != nil {
		writeSSE(writer, "error", errorResponse{Error: err.Error()})
		return
	}
	writeSSE(writer, "report", report)
}

func decodeScenario(writer http.ResponseWriter, request *http.Request) (sim.Scenario, bool) {
	defer request.Body.Close()
	request.Body = http.MaxBytesReader(writer, request.Body, 4<<20)
	var scenario sim.Scenario
	if err := decodeJSON(request, &scenario); err != nil {
		writeJSON(writer, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return sim.Scenario{}, false
	}
	return scenario, true
}

func decodeJSON(request *http.Request, target any) error {
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		if errors.As(err, new(*http.MaxBytesError)) {
			return errors.New("request body too large")
		}
		return err
	}
	return nil
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}

func writeSSE(writer http.ResponseWriter, event string, value any) {
	payload, err := json.Marshal(value)
	if err != nil {
		return
	}
	_, _ = writer.Write([]byte("event: " + event + "\n"))
	_, _ = writer.Write([]byte("id: " + strconv.FormatInt(time.Now().UnixNano(), 10) + "\n"))
	_, _ = writer.Write([]byte("data: " + string(payload) + "\n\n"))
	if flusher, ok := writer.(http.Flusher); ok {
		flusher.Flush()
	}
}
