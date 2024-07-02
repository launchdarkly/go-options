test:
	go clean -i .
	go generate .
	go install .
	go generate ./...
	$(MAKE) lint
	go test -tags testing ./...
	diff test/config_options.go test/golden/config_options.go.txt
	diff test/configWithNoError_options.go test/golden/configWithNoError_options.go.txt
	diff test/configWithBuild_options.go test/golden/configWithBuild_options.go.txt

generate:
	go generate .

lint:
	SKIP=no-commit-to-branch pre-commit run -a

golden:
	mkdir -p test/golden
	cp test/config_options.go test/golden/config_options.go.txt

.PHONY: lint golden test

