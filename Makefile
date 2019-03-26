test:
	go install .
	go generate ./...
	$(MAKE) lint
	go test ./...

lint:
	pre-commit run -a

.PHONY: lint test

