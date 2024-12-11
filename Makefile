# Variables
PORT ?= 8080

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building the application..."
	@go build -o travel-planner main.go

# Run the application
.PHONY: run
run: build
	@./travel-planner

# Install dependencies
.PHONY: make-deps
make-deps:
	@echo "Installing dependencies..."
	@go mod tidy


# Test the application by prompting the user for input
.PHONY: mock-user
mock-user:
	@echo "Testing the /plan endpoint..."
	@read -p "Enter City: " CITY; \
	read -p "Enter Start Date (YYYY-MM-DD): " START_DATE; \
	read -p "Enter End Date (YYYY-MM-DD): " END_DATE; \
	read -p "Enter Minimum Budget: " MIN_BUDGET; \
	read -p "Enter Maximum Budget: " MAX_BUDGET; \
	CURL_CMD="curl \"http://localhost:$(PORT)/plan?city=$$CITY&start-date=$$START_DATE&end-date=$$END_DATE&min-budget=$$MIN_BUDGET&max-budget=$$MAX_BUDGET\""; \
	echo "Executing: $$CURL_CMD"; \
	eval $$CURL_CMD

# Clean the build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -f travel-planner
	@rm -f telemetry.logs
