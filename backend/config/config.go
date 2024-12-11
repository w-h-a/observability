package config

import "os"

var (
    Namespace = os.Getenv("NAMESPACE")
    Name = os.Getenv("NAME")
    Version = os.Getenv("VERSION")
    HttpAddress = os.Getenv("HTTP_ADDRESS")
)