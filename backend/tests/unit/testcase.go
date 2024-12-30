package unit

import "github.com/w-h-a/trace-blame/backend/src/clients/repos"

type TestCase struct {
	When            string
	Endpoint        string
	Query           string
	Client          repos.Client
	Then            string
	ReadCalledTimes int
	ReadCalledWith  []map[string]interface{}
	Payload         string
}
