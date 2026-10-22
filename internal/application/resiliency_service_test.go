package application

import (
	"math/rand"
	"testing"
)

func TestSelectStatusCodeUsesOccurrencesAsWeights(t *testing.T) {
	statusCodes := []uint32{0, 0, 0, 0, 0, 0, 0, 2, 2, 2}
	random := rand.New(rand.NewSource(1))
	counts := map[uint32]int{}

	for range 10_000 {
		code := selectStatusCode(statusCodes, random.Intn)
		counts[code]++
	}

	okRatio := float64(counts[0]) / 10_000
	if okRatio < 0.68 || okRatio > 0.72 {
		t.Fatalf("OK ratio = %.2f, want approximately 0.70", okRatio)
	}
}
