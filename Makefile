GOCMD=go
GOBUILD=$(GOCMD) build

BINARY_NAME=demeda
BUILD_DIR=./build
MAIN_PACKAGE=.

# Default target
.PHONY:
all: build

# Build for current platform
build:
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

# Clean build files
clean:
	rm -rf $(BUILD_DIR)

.PHONY: all build clean