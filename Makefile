GOCMD=go
GOBUILD=$(GOCMD) build

BINARY_NAME=demeda
BUILD_DIR=./build
MAIN_PACKAGE=.

# Default target
.PHONY:
all: clean build run

# Run executable file
run: build
	$(BUILD_DIR)/$(BINARY_NAME)

# Build for current platform
build:
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

# Clean build files
clean:
	@if [ -d "$(BUILD_DIR)" ]; then \
		echo "Удаляю $(BUILD_DIR)"; \
		rm -rf $(BUILD_DIR); \
	fi
	@if [ -f "clinic.db" ]; then \
		echo "Удаляю clinic.db"; \
		rm -rf clinic.db; \
	fi

.PHONY: all build clean