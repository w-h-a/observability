package unit

import "github.com/w-h-a/trace-blame/backend/src/clients/traces"

type TestCase struct {
	When            string
	Endpoint        string
	Query           string
	Client          traces.Client
	Then            string
	ReadCalledTimes int
	ReadCalledWith  []map[string]interface{}
	Payload         string
}
