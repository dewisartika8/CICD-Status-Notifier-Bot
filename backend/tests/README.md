# Test Organization Guide

This document explains the structure and organization of tests in the CICD Status Notifier Bot project.

## ✅ Current Status

**All tests are passing!** 🎉

- **133 tests passing**
- **8 tests skipped** (SQLite tests requiring CGO)
- **Well-organized structure** aligned with source code
- **Centralized mocks** for consistency
- **Test fixtures** for reusable test data

## Directory Structure

```
tests/
├── README.md                 # This file - test organization guide
├── REORGANIZATION_SUMMARY.md # Detailed summary of changes made
├── Makefile                  # Test automation commands
├── fixtures/                 # Test data and helpers
│   └── project_fixtures.go   # Common project test data
├── mocks/                    # Centralized mock implementations
│   └── project_service.go    # Comprehensive project service mock
├── testutils/                # Shared test utilities and helpers
├── integration/              # End-to-end and integration tests
│   ├── handlers/             # HTTP handlers integration tests
│   │   ├── telegram_webhook_integration_test.go
│   │   └── webhook_integration_test.go
│   └── services/             # Service integration tests
│       ├── delivery_integration_test.go
│       └── project_integration_test.go
└── unit/                     # Unit tests organized by domain
    ├── core/                 # Core business logic tests
    │   ├── bot/              # Bot domain tests
    │   │   ├── domain/       # Command validation, routing
    │   │   └── service/      # Bot service logic
    │   ├── notification/     # Notification domain tests
    │   │   ├── delivery/     # Queue, rate limiting, delivery
    │   │   ├── domain/       # Templates, logs, retry config
    │   │   └── service/      # Formatting, retry services
    │   └── project/          # Project domain tests
    │       └── project_service_test.go
    └── adapter/              # Adapter layer tests
        ├── handlers/         # HTTP handlers unit tests (reserved)
        ├── repositories/     # Repository unit tests
        │   └── project_repository_test.go
        └── telegram/         # Telegram adapter tests
            ├── status_command_test.go
            └── status_command_extended_test.go
```

## Test Categories

### 1. Unit Tests (`/unit`)
- Test individual components in isolation
- Use mocks for dependencies
- Fast execution
- High code coverage focus

### 2. Integration Tests (`/integration`)
- Test component interactions
- May use real dependencies (database, external services)
- Slower execution
- Focus on contract testing

### 3. Mocks (`/mocks`)
- Centralized mock implementations
- Shared across unit and integration tests
- Generated and manual mocks

### 4. Test Utils (`/testutils`)
- Common test setup and teardown
- Test data builders
- Shared assertions and helpers

## Naming Conventions

### Test Files
- Unit tests: `*_test.go`
- Integration tests: `*_integration_test.go`
- Mock files: `*_mock.go`

### Test Functions
- Unit tests: `TestComponentName_MethodName_Scenario`
- Integration tests: `TestIntegration_ComponentName_Scenario`

## Running Tests

```bash
# Run all tests
make test

# Run only unit tests
make test-unit

# Run only integration tests
make test-integration

# Run tests with coverage
make test-coverage
```

## Best Practices

1. **Isolation**: Unit tests should not depend on external resources
2. **Clarity**: Test names should clearly describe what is being tested
3. **Arrangement**: Use Given-When-Then or Arrange-Act-Assert patterns
4. **Maintainability**: Keep tests simple and focused on single behaviors
5. **Coverage**: Aim for high code coverage but focus on meaningful tests
