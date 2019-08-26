test:
	go clean -i .
	go generate .
	go install .
	go generate ./...
	$(MAKE) lint
	go test ./...
	diff test/config_options.go test/golden/config_options.go.txt

generate:
	go generate .

lint:
	SKIP=no-commit-to-branch pre-commit run -a

golden:
	mkdir -p test/golden
	cp test/config_options.go test/golden/config_options.go.txt

.PHONY: lint golden test

