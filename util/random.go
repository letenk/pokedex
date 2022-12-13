package util

import (
	"math/rand"
	"strings"
	"time"
)

// Const alphabet for use random data with string
const alphabet = "abcdefghijklmnopqrstuvwxyz"

// Func init for first run
func init() {
	// Run rand.Seed
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	// Get total character on const alphabet
	k := len(alphabet)

	// Loop through n
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomStringFromSet(a ...string) string {
	n := len(a)
	if n == 0 {
		return ""
	}
	return a[rand.Intn(n)]
}
