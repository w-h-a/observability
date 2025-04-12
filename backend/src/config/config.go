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
	namespace          string
	name               string
	version            string
	httpAddress        string
	tracesStore        string
	tracesStoreAddress string
	tracesDB           string
	tracesTable        string
}

func NewConfig() {
	once.Do(func() {
		instance = &config{
			namespace:          "test",
			name:               "test",
			version:            "0.1.0-alpha.0",
			httpAddress:        ":0",
			tracesStore:        "",
			tracesStoreAddress: "",
			tracesDB:           "",
			tracesTable:        "",
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

		store := os.Getenv("TRACES_STORE")
		if len(store) > 0 {
			instance.tracesStore = store
		}

		storeAddress := os.Getenv("TRACES_STORE_ADDRESS")
		if len(storeAddress) > 0 {
			instance.tracesStoreAddress = storeAddress
		}

		db := os.Getenv("TRACES_DB")
		if len(db) > 0 {
			instance.tracesDB = db
		}

		table := os.Getenv("TRACES_TABLE")
		if len(table) > 0 {
			instance.tracesTable = table
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
