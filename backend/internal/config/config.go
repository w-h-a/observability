package config

import (
	"os"
	"sync"

	"github.com/w-h-a/pkg/telemetry/log"
)

var (
	instance *config
	once     sync.Once
)

type config struct {
	namespace           string
	name                string
	version             string
	httpAddress         string
	tracesStore         string
	tracesStoreAddress  string
	tracesDB            string
	tracesTable         string
	metricsStoreAddress string
}

func New() {
	once.Do(func() {
		instance = &config{
			namespace:           "test",
			name:                "test",
			version:             "0.1.0-alpha.0",
			httpAddress:         ":0",
			tracesStore:         "",
			tracesStoreAddress:  "",
			tracesDB:            "",
			tracesTable:         "",
			metricsStoreAddress: "",
		}

		namespace := os.Getenv("NAMESPACE")
		if len(namespace) > 0 {
			instance.namespace = namespace
		}

		name := os.Getenv("NAME")
		if len(name) > 0 {
			instance.name = name
		}

		version := os.Getenv("VERSION")
		if len(version) > 0 {
			instance.version = version
		}

		httpAddress := os.Getenv("HTTP_ADDRESS")
		if len(httpAddress) > 0 {
			instance.httpAddress = httpAddress
		}

		tracesStore := os.Getenv("TRACES_STORE")
		if len(tracesStore) > 0 {
			instance.tracesStore = tracesStore
		}

		tracesStoreAddress := os.Getenv("TRACES_STORE_ADDRESS")
		if len(tracesStoreAddress) > 0 {
			instance.tracesStoreAddress = tracesStoreAddress
		}

		tracesDB := os.Getenv("TRACES_DB")
		if len(tracesDB) > 0 {
			instance.tracesDB = tracesDB
		}

		tracesTable := os.Getenv("TRACES_TABLE")
		if len(tracesTable) > 0 {
			instance.tracesTable = tracesTable
		}

		metricsStoreAddress := os.Getenv("METRICS_STORE_ADDRESS")
		if len(metricsStoreAddress) > 0 {
			instance.metricsStoreAddress = metricsStoreAddress
		}
	})
}

func Namespace() string {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.namespace
}

func Name() string {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.name
}

func Version() string {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.version
}

func HttpAddress() string {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.httpAddress
}

func TracesStore() string {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.tracesStore
}

func TracesStoreAddress() string {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.tracesStoreAddress
}

func TracesDB() string {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.tracesDB
}

func TracesTable() string {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.tracesTable
}

func MetricsStoreAddress() string {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.metricsStoreAddress
}
