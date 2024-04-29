.PHONY: test
test:
	go test ./...

.PHONY: example
example:
	go run ./example/main.go
