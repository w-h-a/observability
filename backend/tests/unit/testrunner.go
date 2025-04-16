package unit

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/observability/backend/internal"
	"github.com/w-h-a/observability/backend/internal/clients/traces/mock"
	"github.com/w-h-a/observability/backend/internal/config"
	"github.com/w-h-a/pkg/telemetry/log"
	memorylog "github.com/w-h-a/pkg/telemetry/log/memory"
	"github.com/w-h-a/pkg/utils/httputils"
	"github.com/w-h-a/pkg/utils/memoryutils"
)

func RunTestCases(t *testing.T, testCases []TestCase) {
	config.New()

	// name
	name := fmt.Sprintf("%s.%s", config.Namespace(), config.Name())

	// log
	logBuffer := memoryutils.NewBuffer()

	logger := memorylog.NewLog(
		log.LogWithPrefix(name),
		memorylog.LogWithBuffer(logBuffer),
	)

	log.SetLogger(logger)

	// traces

	for _, testCase := range testCases {
		var bs []byte
		var err error

		httpServer := internal.Factory(testCase.TracesClient, testCase.MetricsClient)

		mockTracesClient := testCase.TracesClient.(*mock.MockTracesClient)

		t.Run(testCase.When, func(t *testing.T) {
			err = httpServer.Run()
			require.NoError(t, err)

			bs, err = httputils.HttpGet(fmt.Sprintf("http://%s%s%s", httpServer.Options().Address, testCase.Endpoint, testCase.Query))
			require.NoError(t, err)

			t.Log(testCase.Then)

			calls := mockTracesClient.ReadCalledWith()

			require.Equal(t, testCase.ReadCalledTimes, len(calls))

			for i, call := range calls {
				require.Equal(t, testCase.ReadCalledWith[i], call)
			}

			require.Equal(t, testCase.Payload, string(bs))

			t.Cleanup(func() {
				err = httpServer.Stop()
				require.NoError(t, err)

				mockTracesClient.ResetCalledWith()
			})
		})
	}
}
