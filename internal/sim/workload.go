package sim

import (
	"fmt"
	"math"
	"sort"
)

func GenerateWorkload(workload Workload) ([]Assignment, error) {
	if workload.KeyCount <= 0 {
		return nil, fmt.Errorf("keyCount must be positive")
	}
	if workload.Requests <= 0 {
		return nil, fmt.Errorf("requests must be positive")
	}

	prefix := workload.KeyPrefix
	if prefix == "" {
		prefix = "key"
	}

	counts := make(map[string]int, workload.KeyCount)
	remaining := workload.Requests

	for _, hot := range workload.HotKeys {
		if hot.Key == "" {
			return nil, fmt.Errorf("hot key must not be empty")
		}
		if hot.Share < 0 || hot.Share > 1 {
			return nil, fmt.Errorf("hot key share must be between 0 and 1")
		}
		count := hot.Weight
		if count <= 0 {
			count = int(math.Round(float64(workload.Requests) * hot.Share))
		}
		if count > remaining {
			count = remaining
		}
		counts[hot.Key] += count
		remaining -= count
	}

	baseKeys := workload.KeyCount
	if remaining < baseKeys {
		baseKeys = remaining
	}
	for i := range baseKeys {
		key := fmt.Sprintf("%s-%06d", prefix, i)
		counts[key]++
		remaining--
	}

	for remaining > 0 {
		index := hash64(fmt.Sprintf("%d-%d", workload.Seed, remaining)) % uint64(workload.KeyCount)
		key := fmt.Sprintf("%s-%06d", prefix, index)
		counts[key]++
		remaining--
	}

	assignments := make([]Assignment, 0, len(counts))
	for key, count := range counts {
		if count > 0 {
			assignments = append(assignments, Assignment{Key: key, Count: count})
		}
	}
	sort.Slice(assignments, func(i, j int) bool {
		if assignments[i].Count != assignments[j].Count {
			return assignments[i].Count > assignments[j].Count
		}
		return assignments[i].Key < assignments[j].Key
	})
	return assignments, nil
}
