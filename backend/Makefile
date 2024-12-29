.PHONY: tidy
tidy:
	go mod tidy

.PHONY: style
style:
	goimports -l -w ./

.PHONY: unit-test
unit-test:
	go clean -testcache && go test -v ./tests/unit/*/

.PHONY: integration-test
integration-test:
	go clean -testcache && go test -v ./tests/integration/*/

.PHONY: e2e-test
e2e-test:
	go clean -testcache && go test -v ./tests/e2e/*/

.PHONY: go-build
go-build:
	CGO_ENABLED=0 go build -o ./bin/trace-blame-backend ./

.PHONY: go-install
go-install:
	go install

.PHONY: build
build:
	docker buildx build --platform linux/amd64 -t github.com/w-h-a/trace-blame/backend:0.1.1-alpha .
