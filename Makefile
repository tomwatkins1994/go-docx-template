TEST_DIR ?= ./...
COVERAGE_FILE = coverage.out
COVERAGE_HTML = coverage.html

test:
	@echo "Running tests with coverage..."
	@go test -coverprofile=$(COVERAGE_FILE) $(TEST_DIR)

coverage: test
	@echo "Generating HTML coverage report..."
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Opening coverage report in browser..."
	@xdg-open $(COVERAGE_HTML) 2>/dev/null || open $(COVERAGE_HTML) 2>/dev/null || start $(COVERAGE_HTML)

clean:
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@echo "Cleaned up coverage files."
