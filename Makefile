TEST_DIR ?= ./...
COVERAGE_FILE = coverage.out
COVERAGE_HTML = coverage.html

test:
	@echo "Running tests..."
	@go test $(TEST_DIR)

test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=$(COVERAGE_FILE) $(TEST_DIR)
	@echo "Generating HTML coverage report..."
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage file generated at: $(CURDIR)/$(COVERAGE_HTML)"
	@if grep -qi microsoft /proc/version; then \
		powershell.exe Start-Process \"chrome\" \"$(shell wslpath -w $(CURDIR)/$(COVERAGE_HTML) | sed 's/\\/\\\\/g')\" || \
		powershell.exe Start-Process \"msedge\" \"$(shell wslpath -w $(CURDIR)/$(COVERAGE_HTML) | sed 's/\\/\\\\/g')\" || \
		powershell.exe Start-Process \"firefox\" \"$(shell wslpath -w $(CURDIR)/$(COVERAGE_HTML) | sed 's/\\/\\\\/g')\"; \
	elif command -v xdg-open > /dev/null; then \
		xdg-open $(COVERAGE_HTML) > /dev/null 2>&1 & \
	elif command -v open > /dev/null; then \
		open $(COVERAGE_HTML); \
	else \
		echo "Could not detect how to open the coverage report."; \
	fi
	
test-clean:
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@echo "Cleaned up coverage files."

benchmark:
	@go test -bench=Benchmark* -benchtime=1x
