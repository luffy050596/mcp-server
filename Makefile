GOCMD=GO111MODULE=on CGO_ENABLED=0 go
GOBUILD=${GOCMD} build
GOINSTALL=${GOCMD} install
TOOLS_SHELL="./hack/tools.sh"
# golangci-lint
LINTER := bin/golangci-lint

$(LINTER):
	curl -SL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v2.1.1

.PHONY: build
# Build executable file
# Usage: make build dir=<directory_name>
build:
	@if [ -z "$(dir)" ]; then \
		echo "Please specify a directory name using dir=<directory_name>"; \
		exit 1; \
	fi
	@if [ ! -d "$(dir)" ]; then \
		echo "Directory '$(dir)' does not exist"; \
		exit 1; \
	fi
	@mkdir -p bin/
	@echo "Building $(dir)..."
	@$(GOBUILD) -o ./bin/mcp-$(dir) ./$(dir)

.PHONY: build-all
# Build all executables
build-all:
	@rm -rf bin/ && mkdir -p bin/
	@$(GOBUILD) -o ./bin/ ./...

.PHONY: install
# go install
install: build
	@cp bin/mcp-$(dir) $(GOPATH)/bin/

.PHONY: install-all
# go install all
install-all:
	@go mod tidy && $(GOINSTALL)

.PHONY: test
# Test
test:
	go test -v ./... -cover

.PHONY: test-coverage
test-coverage:
	@${TOOLS_SHELL} test_coverage

.PHONY: lint
lint: $(LINTER)
	echo $(os)
	@${TOOLS_SHELL} lint

.PHONY: changelog
# generate changelog
changelog:
	git-chglog -o ./CHANGELOG.md
