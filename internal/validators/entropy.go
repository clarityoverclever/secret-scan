package validators

import "math"

func EntropyValidator(entropyThreshold float64) ValidatorFunc {
	return func(match string, context string) bool {
		return calculateEntropy(match) >= entropyThreshold
	}
}

func calculateEntropy(s string) float64 {
	if len(s) == 0 {
		return 0
	}

	freq := make(map[rune]int)
	for _, c := range s {
		freq[c]++
	}

	var entropy float64
	length := float64(len(s))
	for _, count := range freq {
		p := float64(count) / length
		entropy -= p * math.Log2(p)
	}

	return entropy
}
