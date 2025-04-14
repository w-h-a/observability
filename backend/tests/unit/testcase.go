package unit

import (
	"github.com/w-h-a/observability/backend/internal/clients/metrics"
	"github.com/w-h-a/observability/backend/internal/clients/traces"
)

type TestCase struct {
	When            string
	Endpoint        string
	Query           string
	TracesClient    traces.Client
	MetricsClient   metrics.Client
	Then            string
	ReadCalledTimes int
	ReadCalledWith  []map[string]interface{}
	Payload         string
}
