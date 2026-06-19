package sim

import (
	"crypto/sha256"
	"encoding/binary"
)

func hash64(parts ...string) uint64 {
	h := sha256.New()
	for _, part := range parts {
		_, _ = h.Write([]byte(part))
		_, _ = h.Write([]byte{0})
	}
	sum := h.Sum(nil)
	return binary.BigEndian.Uint64(sum[:8])
}

func hashFloat(parts ...string) float64 {
	value := hash64(parts...)
	return float64(value) / float64(^uint64(0))
}

func uint64Bytes(value uint64) []byte {
	var out [8]byte
	binary.LittleEndian.PutUint64(out[:], value)
	return out[:]
}
