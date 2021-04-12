package keygen

import (
	"math/rand"
	"strings"
	"time"
)

const alphabets = "bcdfghjkmnpqrstvwxyzBCDFGHJKLMNPRSTVWXYZ234567890"
const count = len(alphabets)

// Generates a random 6-character token
func GenerateKey() string {
	seed := rand.NewSource(time.Now().UnixNano())
	var str strings.Builder
	str.Grow(6)
	r := rand.New(seed)
	for i := 0; i < 6; i++ {
		str.WriteByte(alphabets[r.Intn(count)])
	}
	return str.String()
}
