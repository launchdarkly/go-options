test:
	go install .
	go generate ./...
	$(MAKE) lint
	go test ./...

lint:
	SKIP=no-commit-to-branch pre-commit run -a

.PHONY: lint test

