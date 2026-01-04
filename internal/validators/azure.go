package validators

import "strings"

var AzureContextValidator = ValidatorFunc(func(match string, context string) bool {
	lower := strings.ToLower(context)
	keywords := []string{
		"azure",
		"azure_sub",
		"azure_subscription",
		"subscription",
		"subscriptionid",
		"tenant",
	}

	for _, keyword := range keywords {
		if strings.Contains(lower, keyword) {
			return true
		}
	}
	return false
})
