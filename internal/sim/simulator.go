package sim

import (
	"fmt"
	"math"
	"sort"
)

func Simulate(scenario Scenario) (Report, error) {
	if scenario.Name == "" {
		scenario.Name = "shardlab-scenario"
	}
	assignments, err := GenerateWorkload(scenario.Workload)
	if err != nil {
		return Report{}, err
	}

	beforeAssigner, err := NewAssigner(scenario.Algorithm, scenario.NodesBefore, scenario.Ring)
	if err != nil {
		return Report{}, fmt.Errorf("before assigner: %w", err)
	}
	afterAssigner, err := NewAssigner(scenario.Algorithm, scenario.NodesAfter, scenario.Ring)
	if err != nil {
		return Report{}, fmt.Errorf("after assigner: %w", err)
	}

	before := assignAll(assignments, beforeAssigner)
	after := assignAll(assignments, afterAssigner)
	moved := movedKeys(before, after)

	totalRequests := 0
	movedRequests := 0
	for _, item := range assignments {
		totalRequests += item.Count
		if before[item.Key].NodeID != after[item.Key].NodeID {
			movedRequests += item.Count
		}
	}

	report := Report{
		Scenario:       scenario.Name,
		Algorithm:      scenario.Algorithm,
		Keys:           len(assignments),
		Requests:       totalRequests,
		MovementRatio:  round4(float64(movedRequests) / float64(totalRequests)),
		Before:         summarize(before, scenario.NodesBefore),
		After:          summarize(after, scenario.NodesAfter),
		TopMovedKeys:   topMoved(moved, 8),
		Recommendation: recommendation(scenario.Algorithm, movedRequests, totalRequests, summarize(after, scenario.NodesAfter)),
	}
	return report, nil
}

func assignAll(assignments []Assignment, assigner Assigner) map[string]Assignment {
	result := make(map[string]Assignment, len(assignments))
	for _, assignment := range assignments {
		nodeID, err := assigner.Assign(assignment.Key)
		if err != nil {
			continue
		}
		result[assignment.Key] = Assignment{
			Key:    assignment.Key,
			NodeID: nodeID,
			Count:  assignment.Count,
		}
	}
	return result
}

func movedKeys(before map[string]Assignment, after map[string]Assignment) []MovedKey {
	var moved []MovedKey
	for key, beforeAssignment := range before {
		afterAssignment, ok := after[key]
		if !ok || beforeAssignment.NodeID == afterAssignment.NodeID {
			continue
		}
		moved = append(moved, MovedKey{
			Key:        key,
			Count:      beforeAssignment.Count,
			BeforeNode: beforeAssignment.NodeID,
			AfterNode:  afterAssignment.NodeID,
		})
	}
	sort.Slice(moved, func(i, j int) bool {
		if moved[i].Count != moved[j].Count {
			return moved[i].Count > moved[j].Count
		}
		return moved[i].Key < moved[j].Key
	})
	return moved
}

func topMoved(moved []MovedKey, limit int) []MovedKey {
	if len(moved) <= limit {
		return moved
	}
	return append([]MovedKey(nil), moved[:limit]...)
}

func summarize(assignments map[string]Assignment, nodes []Node) Distribution {
	loadsByNode := map[string]int{}
	for _, node := range nodes {
		loadsByNode[node.ID] = 0
	}
	total := 0
	for _, assignment := range assignments {
		loadsByNode[assignment.NodeID] += assignment.Count
		total += assignment.Count
	}

	nodeLoads := make([]NodeLoad, 0, len(loadsByNode))
	values := make([]int, 0, len(loadsByNode))
	for nodeID, load := range loadsByNode {
		values = append(values, load)
		share := 0.0
		if total > 0 {
			share = float64(load) / float64(total)
		}
		nodeLoads = append(nodeLoads, NodeLoad{NodeID: nodeID, Load: load, Share: round4(share)})
	}
	sort.Slice(nodeLoads, func(i, j int) bool { return nodeLoads[i].NodeID < nodeLoads[j].NodeID })
	sort.Ints(values)

	mean := 0.0
	if len(values) > 0 {
		mean = float64(total) / float64(len(values))
	}
	minLoad := 0
	maxLoad := 0
	if len(values) > 0 {
		minLoad = values[0]
		maxLoad = values[len(values)-1]
	}
	skew := 0.0
	if mean > 0 {
		skew = float64(maxLoad) / mean
	}

	var overloaded []string
	var underloaded []string
	for _, load := range nodeLoads {
		if mean > 0 && float64(load.Load) > mean*1.25 {
			overloaded = append(overloaded, load.NodeID)
		}
		if mean > 0 && float64(load.Load) < mean*0.75 {
			underloaded = append(underloaded, load.NodeID)
		}
	}

	return Distribution{
		NodeLoads:        nodeLoads,
		MeanLoad:         round2(mean),
		MaxLoad:          maxLoad,
		MinLoad:          minLoad,
		LoadSkew:         round4(skew),
		GiniCoefficient:  round4(gini(values)),
		OverloadedNodes:  overloaded,
		UnderloadedNodes: underloaded,
	}
}

func recommendation(algorithm Algorithm, movedRequests int, totalRequests int, after Distribution) string {
	movement := float64(movedRequests) / float64(totalRequests)
	switch {
	case algorithm == AlgorithmModulo && movement > 0.4:
		return "Modulo hashing moved a large share of requests; use consistent or rendezvous hashing before adding/removing nodes."
	case after.LoadSkew > 1.35:
		return "Load is skewed after reassignment; increase virtual nodes, split hot keys, or use weighted nodes."
	case movement < 0.3 && after.LoadSkew <= 1.25:
		return "Reassignment looks healthy: limited key movement and balanced node load."
	default:
		return "Review hot keys and node weights before applying this topology change."
	}
}

func gini(values []int) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0
	for _, value := range values {
		total += value
	}
	if total == 0 {
		return 0
	}
	var cumulativeDiff int
	for _, left := range values {
		for _, right := range values {
			cumulativeDiff += int(math.Abs(float64(left - right)))
		}
	}
	return float64(cumulativeDiff) / (2 * float64(len(values)) * float64(total))
}

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}

func round4(value float64) float64 {
	return math.Round(value*10000) / 10000
}
