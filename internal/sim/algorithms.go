package sim

import (
	"fmt"
	"sort"
)

type Assigner interface {
	Assign(key string) (string, error)
}

func NewAssigner(algorithm Algorithm, nodes []Node, ring RingConfig) (Assigner, error) {
	if len(nodes) == 0 {
		return nil, fmt.Errorf("nodes must not be empty")
	}
	for _, node := range nodes {
		if node.ID == "" {
			return nil, fmt.Errorf("node id must not be empty")
		}
	}

	switch algorithm {
	case AlgorithmModulo:
		return moduloAssigner{nodes: normalizeNodes(nodes)}, nil
	case AlgorithmConsistent:
		virtualNodes := ring.VirtualNodes
		if virtualNodes <= 0 {
			virtualNodes = 128
		}
		return newConsistentAssigner(nodes, virtualNodes), nil
	case AlgorithmRendezvous:
		return rendezvousAssigner{nodes: normalizeNodes(nodes)}, nil
	default:
		return nil, fmt.Errorf("unknown algorithm %q", algorithm)
	}
}

type moduloAssigner struct {
	nodes []Node
}

func (m moduloAssigner) Assign(key string) (string, error) {
	if len(m.nodes) == 0 {
		return "", fmt.Errorf("nodes must not be empty")
	}
	index := hash64(key) % uint64(len(m.nodes))
	return m.nodes[index].ID, nil
}

type ringPoint struct {
	hash   uint64
	nodeID string
}

type consistentAssigner struct {
	points []ringPoint
}

func newConsistentAssigner(nodes []Node, virtualNodes int) consistentAssigner {
	points := make([]ringPoint, 0, len(nodes)*virtualNodes)
	for _, node := range normalizeNodes(nodes) {
		weight := node.Weight
		if weight <= 0 {
			weight = 1
		}
		replicas := max(1, int(float64(virtualNodes)*weight))
		for replica := range replicas {
			points = append(points, ringPoint{
				hash:   hash64(node.ID, string(uint64Bytes(uint64(replica)))),
				nodeID: node.ID,
			})
		}
	}
	sort.Slice(points, func(i, j int) bool { return points[i].hash < points[j].hash })
	return consistentAssigner{points: points}
}

func (c consistentAssigner) Assign(key string) (string, error) {
	if len(c.points) == 0 {
		return "", fmt.Errorf("ring has no points")
	}
	value := hash64(key)
	index := sort.Search(len(c.points), func(i int) bool { return c.points[i].hash >= value })
	if index == len(c.points) {
		index = 0
	}
	return c.points[index].nodeID, nil
}

type rendezvousAssigner struct {
	nodes []Node
}

func (r rendezvousAssigner) Assign(key string) (string, error) {
	if len(r.nodes) == 0 {
		return "", fmt.Errorf("nodes must not be empty")
	}
	bestNode := ""
	bestScore := -1.0
	for _, node := range r.nodes {
		weight := node.Weight
		if weight <= 0 {
			weight = 1
		}
		score := hashFloat(key, node.ID) * weight
		if score > bestScore {
			bestScore = score
			bestNode = node.ID
		}
	}
	return bestNode, nil
}

func normalizeNodes(nodes []Node) []Node {
	out := append([]Node(nil), nodes...)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
