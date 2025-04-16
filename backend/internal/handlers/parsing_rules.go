package handlers

var (
	SupportedRequestedDimensions = []string{"calls", "duration"}

	SupportedRequestedAggregations = map[string][]string{
		"calls":    {"count", "rate_per_sec"},
		"duration": {"avg", "p50", "p95", "p99"},
	}

	SupportedRequestedTagOperators = []string{"equals", "contains", "isnotnull"}

	SupportedMetrics = []string{"cpu"}

	RequiredMetricAttributes = map[string][]string{
		"cpu": {"service"},
	}
)
