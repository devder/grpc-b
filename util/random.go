package util

import (
	"math"
	"math/rand"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var rng *rand.Rand

// A random number generator, seeded with the current time (time.Now().UnixNano()),
// to ensure different sequences each time the function is called.
func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int) int64 {
	// Intn returns a non-negative pseudo-random number in [0, n)
	return int64(rng.Intn(max-min) + min)
}

func RandomFloat64(min, max float64) float64 {
	// Float64 returns, as a float64, a pseudo-random number in [0.0,1.0)

	result := min + rng.Float64()*(max-min)
	return math.Round(result*10_000) / 10_000 // to 4dp
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	b := make([]byte, n)

	for i := range b {
		b[i] = alphabet[rng.Intn(len(alphabet))]
	}

	return string(b)
}

// Generate a random owner name
func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() float64 {
	return float64(RandomFloat64(50.55, 999.99))
}

func RandomCurrency() string {
	currencies := [3]string{EUR, USD, CAD}

	return currencies[rng.Intn(len(currencies))]
}
