#!/bin/bash
set -e

# Main integration test orchestrator for camp
# Builds Docker image, compiles camp binary, and runs all test scenarios

# Color output for better readability
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
DOCKER_DIR="$PROJECT_ROOT/test/integration/docker"
FIXTURES_DIR="$PROJECT_ROOT/test/integration/fixtures"

# Docker image name
IMAGE_NAME="camp-integration-test:latest"

# Test results tracking
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0
declare -a FAILED_TESTS

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}Camp Integration Test Suite${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""

# Step 1: Build camp binary for Linux (cross-compile if on macOS)
echo -e "${YELLOW}[1/3] Building camp binary for Linux...${NC}"
cd "$PROJECT_ROOT"
if CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o camp main.go; then
    echo -e "${GREEN}✓ Camp binary built successfully (linux/amd64)${NC}"
else
    echo -e "${RED}✗ Failed to build camp binary${NC}"
    exit 1
fi
echo ""

# Step 2: Build Docker image
echo -e "${YELLOW}[2/3] Building Docker image...${NC}"
cd "$DOCKER_DIR"
if docker build -t "$IMAGE_NAME" .; then
    echo -e "${GREEN}✓ Docker image built successfully${NC}"
else
    echo -e "${RED}✗ Failed to build Docker image${NC}"
    exit 1
fi
echo ""

# Step 3: Run test scripts
echo -e "${YELLOW}[3/3] Running integration tests...${NC}"
echo ""

# Function to run a single test script
run_test() {
    local test_name="$1"
    local test_script="$2"

    echo -e "${BLUE}▶ Running: $test_name${NC}"
    ((TESTS_RUN++))

    # Run test in fresh container
    # Mount camp binary, test scripts, fixtures, and templates
    if docker run --rm \
        -v "$PROJECT_ROOT/camp:/home/testuser/bin/camp:ro" \
        -v "$SCRIPT_DIR:/home/testuser/tests:ro" \
        -v "$FIXTURES_DIR:/home/testuser/fixtures:ro" \
        -v "$PROJECT_ROOT/templates:/home/testuser/templates:ro" \
        "$IMAGE_NAME" \
        /bin/bash /home/testuser/tests/"$test_script"; then

        echo -e "${GREEN}✓ PASSED: $test_name${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}✗ FAILED: $test_name${NC}"
        ((TESTS_FAILED++))
        FAILED_TESTS+=("$test_name")
    fi
    echo ""
}

# Run all test scenarios
run_test "Bootstrap" "test-bootstrap.sh"
run_test "Environment Rebuild" "test-rebuild.sh"
run_test "Package Installation" "test-packages.sh"
run_test "Flake Integration" "test-flakes.sh"
run_test "Environment Cleanup" "test-nuke.sh"

# Print summary
echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}Test Summary${NC}"
echo -e "${BLUE}======================================${NC}"
echo -e "Total Tests:  $TESTS_RUN"
echo -e "${GREEN}Passed:       $TESTS_PASSED${NC}"
echo -e "${RED}Failed:       $TESTS_FAILED${NC}"

if [ ${#FAILED_TESTS[@]} -gt 0 ]; then
    echo ""
    echo -e "${RED}Failed Tests:${NC}"
    for test in "${FAILED_TESTS[@]}"; do
        echo -e "${RED}  - $test${NC}"
    done
fi

echo -e "${BLUE}======================================${NC}"

# Exit with appropriate code
if [ $TESTS_FAILED -gt 0 ]; then
    exit 1
else
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
fi
