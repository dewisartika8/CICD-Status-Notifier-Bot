#!/usr/bin/env bash

echo "🚀 Running CICD Status Notifier Bot Tests"
echo "========================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Function to run tests and track results
run_test_suite() {
    local test_path=$1
    local test_name=$2
    
    echo -e "\n${YELLOW}📋 Running $test_name...${NC}"
    
    if go test "$test_path" -v; then
        echo -e "${GREEN}✅ $test_name PASSED${NC}"
        ((PASSED_TESTS++))
    else
        echo -e "${RED}❌ $test_name FAILED${NC}"
        ((FAILED_TESTS++))
    fi
    ((TOTAL_TESTS++))
}

# Change to backend directory
cd "$(dirname "$0")" || exit 1

echo "📁 Current directory: $(pwd)"
echo ""

# Run all test suites
echo "🧪 Starting Test Execution..."

# Unit Tests - Bot Domain
run_test_suite "./tests/unit/bot/domain/..." "Bot Domain Layer Tests"

# Unit Tests - Bot Service  
run_test_suite "./tests/unit/bot/service/..." "Bot Service Layer Tests"

# Integration Tests - Bot
run_test_suite "./tests/integration/bot/..." "Bot Integration Tests"

# Unit Tests - Core (existing)
run_test_suite "./tests/unit/..." "Core Unit Tests"

# Integration Tests - Webhook (existing)
run_test_suite "./tests/integration/..." "Webhook Integration Tests"

# Summary
echo ""
echo "📊 TEST SUMMARY"
echo "==============="
echo -e "Total Test Suites: $TOTAL_TESTS"
echo -e "${GREEN}Passed: $PASSED_TESTS${NC}"

if [ $FAILED_TESTS -gt 0 ]; then
    echo -e "${RED}Failed: $FAILED_TESTS${NC}"
    echo ""
    echo -e "${RED}❌ Some tests failed. Please check the output above.${NC}"
    exit 1
else
    echo -e "${GREEN}Failed: $FAILED_TESTS${NC}"
    echo ""
    echo -e "${GREEN}🎉 All tests passed successfully!${NC}"
    echo ""
    echo "✨ Bot implementation is ready for deployment!"
fi
