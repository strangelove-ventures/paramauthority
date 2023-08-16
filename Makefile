.PHONY: proto-all proto-format proto-lint proto-gen format lint test build
all: proto-all format lint build

###############################################################################
###                                  Build                                  ###
###############################################################################

BUILD_FLAGS := -ldflags '$(strip $(LDFLAGS))' -mod readonly -tags 'ledger' -trimpath

build:
	@echo "🤖 Building simd..."
	@cd simapp && go build $(BUILD_FLAGS) -o "$(PWD)/build/" ./cmd/simd && cd ..
	@echo "✅ Completed build!"

install:
	@echo "🤖 Installing simd..."
	@cd simapp && go install $(BUILD_FLAGS) ./cmd/simd && cd ..
	@echo "✅ Completed install!"

###############################################################################
###                          Formatting & Linting                           ###
###############################################################################

gofumpt_cmd=mvdan.cc/gofumpt
golangci_lint_cmd=github.com/golangci/golangci-lint/cmd/golangci-lint

format:
	@echo "🤖 Running formatter..."
	@go run $(gofumpt_cmd) -l -w .
	@echo "✅ Completed formatting!"

lint:
	@echo "🤖 Running linter..."
	@go run $(golangci_lint_cmd) run --timeout=10m
	@echo "✅ Completed linting!"

###############################################################################
###                                Protobuf                                 ###
###############################################################################

BUF_VERSION=1.26.1
BUILDER_VERSION=0.13.5

proto-all: proto-format proto-lint proto-gen

proto-format:
	@echo "🤖 Running protobuf formatter..."
	@docker run --rm --volume "$(PWD)":/workspace --workdir /workspace \
		bufbuild/buf:$(BUF_VERSION) format --diff --write
	@echo "✅ Completed protobuf formatting!"

proto-gen:
	@echo "🤖 Generating code from protobuf..."
	@docker run --rm --volume "$(PWD)":/workspace --workdir /workspace \
		ghcr.io/cosmos/proto-builder:$(BUILDER_VERSION) sh ./proto/generate.sh
	@echo "✅ Completed code generation!"

proto-lint:
	@echo "🤖 Running protobuf linter..."
	@docker run --rm --volume "$(PWD)":/workspace --workdir /workspace \
		bufbuild/buf:$(BUF_VERSION) lint
	@echo "✅ Completed protobuf linting!"

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

heighliner:
	@echo "🤖 Building simd image..."
	@heighliner build --chain authority-simd --file "$(PWD)/chains.yaml" --local
	@echo "✅ Completed build!"

test:
	@echo "🤖 Running tests..."
	@cd tests && go test -mod=readonly ./...
	@echo "✅ Completed tests!"
