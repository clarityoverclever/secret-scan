package validators

import "encoding/base64"

func Base64HighEntropyValidator(entropyThreshold float64) ValidatorFunc {
	return func(match string, context string) bool {
		// Check if valid Base64
		decoded, err := base64.StdEncoding.DecodeString(match)
		if err != nil {
			return false
		}
		// Check if decoded data has high entropy (likely encrypted/random)
		return calculateEntropy(string(decoded)) > entropyThreshold
	}
}
