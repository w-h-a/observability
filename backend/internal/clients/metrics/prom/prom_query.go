package prom

import "fmt"

var (
	Queries = map[string]func(service string) string{
		"cpu": func(service string) string {
			return fmt.Sprintf(`container_cpu_utilization_ratio{container_name="%s"}`, service)
		},
	}
)
