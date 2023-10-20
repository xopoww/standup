package identifiers

import (
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().Unix()))

func GenerateID() string {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], rng.Uint64())
	return hex.EncodeToString(buf[:])
}
