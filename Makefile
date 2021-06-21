ifdef GOBIN
GO_OPTIONS_PATH=$(GOBIN)
else
ifdef GOPATH
GO_OPTIONS_PATH=$(GOPATH)/bin
else
GO_OPTIONS_PATH=$(HOME)/go/bin
endif
endif
	
test: install
	$(MAKE) generate_tests
	$(MAKE) lint
	go test ./...
	diff test/config_options.go test/golden/config_options.go.txt

generate_tests: install
	env PATH="$(GO_OPTIONS_PATH):$(PATH)" go generate -v -x ./test/...

install:
	go clean -i .
	go generate .
	go install .

generate:
	go generate .

lint:
ifdef SKIP_LINT
	echo Skipping lint
else
	SKIP=no-commit-to-branch pre-commit run -a
endif

golden:
	mkdir -p test/golden
	cp test/config_options.go test/golden/config_options.go.txt

.PHONY: install generate generate_tests lint golden test

