# Go parameters
GOCMD=go
GOBUILD=CGO_ENABLED=0 $(GOCMD) build
GOMOD=$(GOCMD) mod
GOTEST=$(GOCMD) test
GOFLAGS := -v 
LDFLAGS := -s -w

APP_NAME := naabu
SRC := ./cmd/naabu
BUILD_DIR := dist

PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	linux/arm \
	linux/386 \
	windows/amd64 \
	windows/arm64 \
	windows/386 \
	darwin/amd64 \
	darwin/arm64 \
	freebsd/amd64 \
	freebsd/arm \
	freebsd/arm64

.PHONY: all clean dist

all: build
build:
	$(GOBUILD) $(GOFLAGS) -ldflags '$(LDFLAGS)' -o "naabu" $(SRC)

dist:
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		output="$(BUILD_DIR)/$(APP_NAME)-$$GOOS-$$GOARCH"; \
		if [ "$$GOOS" = "windows" ]; then output="$$output.exe"; fi; \
		echo "Building $$output"; \
		CGO_ENABLED=0 GOOS=$$GOOS GOARCH=$$GOARCH go build -trimpath -ldflags '$(LDFLAGS)' -tags netgo,osusergo,static_build -o $$output $(SRC); \
	done

test:
	$(GOTEST) $(GOFLAGS) ./...
tidy:
	$(GOMOD) tidy

cleandist:
	rm -rf $(BUILD_DIR)

