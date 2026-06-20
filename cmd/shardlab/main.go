package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kamilch1k/shardlab/internal/sim"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("shardlab", flag.ContinueOnError)
	flags.SetOutput(stderr)

	input := flags.String("input", "", "scenario JSON file")
	output := flags.String("out", "", "optional report JSON path")
	format := flags.String("format", "text", "output format: text or json")
	if err := flags.Parse(args); err != nil {
		return 64
	}
	if strings.TrimSpace(*input) == "" {
		_, _ = fmt.Fprintln(stderr, "missing required -input")
		return 64
	}

	file, err := os.Open(*input)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "open scenario: %v\n", err)
		return 65
	}
	defer file.Close()

	var scenario sim.Scenario
	if err := json.NewDecoder(file).Decode(&scenario); err != nil {
		_, _ = fmt.Fprintf(stderr, "decode scenario: %v\n", err)
		return 65
	}
	report, err := sim.Simulate(scenario)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "simulate: %v\n", err)
		return 65
	}

	if strings.TrimSpace(*output) != "" {
		if err := writeJSON(*output, report); err != nil {
			_, _ = fmt.Fprintf(stderr, "write report: %v\n", err)
			return 65
		}
	}

	switch strings.ToLower(strings.TrimSpace(*format)) {
	case "json":
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(report); err != nil {
			_, _ = fmt.Fprintf(stderr, "write json: %v\n", err)
			return 65
		}
	case "text":
		writeText(stdout, report)
	default:
		_, _ = fmt.Fprintf(stderr, "unknown -format %q\n", *format)
		return 64
	}

	return 0
}

func writeText(writer io.Writer, report sim.Report) {
	_, _ = fmt.Fprintf(
		writer,
		"scenario=%s algorithm=%s keys=%d requests=%d movement=%.2f%% after_skew=%.3f after_gini=%.3f\n",
		report.Scenario,
		report.Algorithm,
		report.Keys,
		report.Requests,
		report.MovementRatio*100,
		report.After.LoadSkew,
		report.After.GiniCoefficient,
	)
	_, _ = fmt.Fprintf(writer, "recommendation=%s\n", report.Recommendation)
	_, _ = fmt.Fprintln(writer, "after_loads:")
	for _, load := range report.After.NodeLoads {
		_, _ = fmt.Fprintf(writer, "  %-16s load=%d share=%.2f%%\n", load.NodeID, load.Load, load.Share*100)
	}
	if len(report.TopMovedKeys) > 0 {
		_, _ = fmt.Fprintln(writer, "top_moved_keys:")
		for _, moved := range report.TopMovedKeys {
			_, _ = fmt.Fprintf(writer, "  %-18s count=%d %s -> %s\n", moved.Key, moved.Count, moved.BeforeNode, moved.AfterNode)
		}
	}
}

func writeJSON(path string, report sim.Report) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}
