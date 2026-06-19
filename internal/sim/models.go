package sim

type Algorithm string

const (
	AlgorithmModulo     Algorithm = "modulo"
	AlgorithmConsistent Algorithm = "consistent"
	AlgorithmRendezvous Algorithm = "rendezvous"
)

type Scenario struct {
	Name        string     `json:"name"`
	Algorithm   Algorithm  `json:"algorithm"`
	NodesBefore []Node     `json:"nodesBefore"`
	NodesAfter  []Node     `json:"nodesAfter"`
	Workload    Workload   `json:"workload"`
	Ring         RingConfig `json:"ring"`
}

type Node struct {
	ID     string  `json:"id"`
	Weight float64 `json:"weight,omitempty"`
}

type RingConfig struct {
	VirtualNodes int `json:"virtualNodes,omitempty"`
}

type Workload struct {
	KeyCount   int      `json:"keyCount"`
	Requests  int      `json:"requests"`
	HotKeys   []HotKey `json:"hotKeys,omitempty"`
	Seed      uint64   `json:"seed,omitempty"`
	KeyPrefix string   `json:"keyPrefix,omitempty"`
}

type HotKey struct {
	Key    string  `json:"key"`
	Share  float64 `json:"share"`
	Weight int     `json:"weight,omitempty"`
}

type Assignment struct {
	Key    string `json:"key"`
	NodeID string `json:"nodeId"`
	Count  int    `json:"count"`
}

type Report struct {
	Scenario       string       `json:"scenario"`
	Algorithm      Algorithm    `json:"algorithm"`
	Keys           int          `json:"keys"`
	Requests       int          `json:"requests"`
	MovementRatio  float64      `json:"movementRatio"`
	Before         Distribution `json:"before"`
	After          Distribution `json:"after"`
	TopMovedKeys   []MovedKey   `json:"topMovedKeys"`
	Recommendation string       `json:"recommendation"`
}

type Distribution struct {
	NodeLoads         []NodeLoad `json:"nodeLoads"`
	MeanLoad         float64    `json:"meanLoad"`
	MaxLoad          int        `json:"maxLoad"`
	MinLoad          int        `json:"minLoad"`
	LoadSkew         float64    `json:"loadSkew"`
	GiniCoefficient  float64    `json:"giniCoefficient"`
	OverloadedNodes  []string   `json:"overloadedNodes"`
	UnderloadedNodes []string   `json:"underloadedNodes"`
}

type NodeLoad struct {
	NodeID string  `json:"nodeId"`
	Load   int     `json:"load"`
	Share  float64 `json:"share"`
}

type MovedKey struct {
	Key        string `json:"key"`
	Count      int    `json:"count"`
	BeforeNode string `json:"beforeNode"`
	AfterNode  string `json:"afterNode"`
}
