package splitify

import (
	"testing"
)

func TestWeighted_Next(t *testing.T) {
	w := NewWeightedSplit()
	w.AddRule(&Rule{
		Handler: "a",
		Weight:  3,
	})
	w.AddRule(&Rule{
		Handler: "b",
		Weight:  97,
	})

	resp := make([]string, 200)
	for i := 0; i < 200; i++ {
		h, err := w.Next()
		if err != nil {
			t.Fatal(err)
		}
		resp[i] = h.(string)
	}

	// check count of a and b
	aCount := 0
	bCount := 0

	for _, v := range resp {
		if v == "a" {
			aCount++
		} else if v == "b" {
			bCount++
		}
	}

	if bCount/aCount != 32 {
		t.Fatalf("expected ratio of b to a to be 32, got %d", bCount/aCount)
	}
}
