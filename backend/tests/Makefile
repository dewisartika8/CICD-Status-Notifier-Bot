# Test Makefile for CICD Status Notifier Bot
# This file provides convenient commands for running tests

.PHONY: test test-unit test-integration test-coverage test-verbose clean help

# Default target
help:
	@echo "Available commands:"
	@echo "  test           Run all tests"
	@echo "  test-unit      Run only unit tests"
	@echo "  test-integration Run only integration tests"
	@echo "  test-coverage  Run tests with coverage report"
	@echo "  test-verbose   Run tests with verbose output"
	@echo "  clean          Clean test cache and coverage files"

# Run all tests
test:
	@echo "Running all tests..."
	go test ./tests/...

# Run only unit tests
test-unit:
	@echo "Running unit tests..."
	go test ./tests/unit/...

# Run only integration tests  
test-integration:
	@echo "Running integration tests..."
	go test ./tests/integration/...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./tests/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v ./tests/...

# Run specific test package
test-telegram:
	@echo "Running telegram tests..."
	go test -v ./tests/unit/adapter/telegram/...

test-project:
	@echo "Running project tests..."
	go test -v ./tests/unit/core/project/...

test-repositories:
	@echo "Running repository tests..."
	go test -v ./tests/unit/adapter/repositories/...

# Clean test artifacts
clean:
	@echo "Cleaning test artifacts..."
	go clean -testcache
	rm -f coverage.out coverage.html

# Run tests in watch mode (requires entr or similar tool)
test-watch:
	@echo "Running tests in watch mode..."
	find . -name "*.go" | entr -r make test

# Benchmark tests
test-bench:
	@echo "Running benchmark tests..."
	go test -bench=. ./tests/...
