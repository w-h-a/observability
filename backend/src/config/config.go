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
	namespace    string
	name         string
	version      string
	httpAddress  string
	store        string
	storeAddress string
	db           string
	table        string
}

func NewConfig() {
	once.Do(func() {
		instance = &config{
			namespace:    "test",
			name:         "test",
			version:      "0.1.0-alpha.0",
			httpAddress:  ":4000",
			store:        "",
			storeAddress: "",
			db:           "",
			table:        "",
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

		store := os.Getenv("STORE")
		if len(store) > 0 {
			instance.store = store
		}

		storeAddress := os.Getenv("STORE_ADDRESS")
		if len(storeAddress) > 0 {
			instance.storeAddress = storeAddress
		}

		db := os.Getenv("DB")
		if len(db) > 0 {
			instance.db = db
		}

		table := os.Getenv("TABLE")
		if len(table) > 0 {
			instance.table = table
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

func Store() string {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.store
}

func StoreAddress() string {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.storeAddress
}

func DB() string {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.db
}

func Table() string {
	if instance == nil {
		log.Fatal("no config instance")
	}

	return instance.table
}
