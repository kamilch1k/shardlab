package sim

import "testing"

func TestModuloMovesManyKeysWhenNodeAdded(t *testing.T) {
	report, err := Simulate(sampleScenario(AlgorithmModulo))
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	if report.MovementRatio < 0.35 {
		t.Fatalf("expected high movement for modulo hashing, got %.4f", report.MovementRatio)
	}
}

func TestConsistentHashingLimitsMovement(t *testing.T) {
	report, err := Simulate(sampleScenario(AlgorithmConsistent))
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	if report.MovementRatio > 0.35 {
		t.Fatalf("expected lower movement for consistent hashing, got %.4f", report.MovementRatio)
	}
	if len(report.After.NodeLoads) != 4 {
		t.Fatalf("expected four after nodes, got %#v", report.After.NodeLoads)
	}
}

func TestRendezvousHashingProducesBalancedLoads(t *testing.T) {
	report, err := Simulate(sampleScenario(AlgorithmRendezvous))
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	if report.After.LoadSkew > 1.45 {
		t.Fatalf("expected acceptable load skew, got %.4f", report.After.LoadSkew)
	}
}

func TestHotKeysCreateVisibleSkew(t *testing.T) {
	scenario := sampleScenario(AlgorithmConsistent)
	scenario.Workload.HotKeys = []HotKey{{Key: "tenant-vip", Share: 0.4}}
	report, err := Simulate(scenario)
	if err != nil {
		t.Fatalf("simulate: %v", err)
	}
	if report.After.LoadSkew < 1.3 {
		t.Fatalf("expected hot-key load skew, got %.4f", report.After.LoadSkew)
	}
}

func sampleScenario(algorithm Algorithm) Scenario {
	return Scenario{
		Name:      "add-cache-node",
		Algorithm: algorithm,
		NodesBefore: []Node{
			{ID: "cache-a"},
			{ID: "cache-b"},
			{ID: "cache-c"},
		},
		NodesAfter: []Node{
			{ID: "cache-a"},
			{ID: "cache-b"},
			{ID: "cache-c"},
			{ID: "cache-d"},
		},
		Workload: Workload{
			KeyCount:  10_000,
			Requests:  50_000,
			Seed:      42,
			KeyPrefix: "tenant",
		},
		Ring: RingConfig{VirtualNodes: 256},
	}
}
