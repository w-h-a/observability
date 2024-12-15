package config

import "os"

var (
	Namespace    = os.Getenv("NAMESPACE")
	Name         = os.Getenv("NAME")
	Version      = os.Getenv("VERSION")
	HttpAddress  = os.Getenv("HTTP_ADDRESS")
	Store        = os.Getenv("STORE")
	StoreAddress = os.Getenv("STORE_ADDRESS")
	DB           = os.Getenv("DB")
	Table        = os.Getenv("TABLE")
)
