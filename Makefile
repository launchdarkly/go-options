test:
	go clean -i .
	go install .
	go generate ./...
	$(MAKE) lint
	go test ./test
	diff test/config_options.go test/golden/config_options.go

lint:
	SKIP=no-commit-to-branch pre-commit run -a

golden:
	mkdir -p test/golden
	cp test/config_options.go test/golden/

.PHONY: lint golden test

