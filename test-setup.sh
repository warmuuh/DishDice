#!/bin/bash

# DishDice Setup Test Script
# This script verifies your local setup is working correctly

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_header() {
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

run_test() {
    local test_name="$1"
    local test_command="$2"

    TOTAL_TESTS=$((TOTAL_TESTS + 1))

    if eval "$test_command" > /dev/null 2>&1; then
        print_success "$test_name"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        return 0
    else
        print_error "$test_name"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        return 1
    fi
}

echo ""
echo -e "${BLUE}╔═══════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                                           ║${NC}"
echo -e "${BLUE}║     🎲 DishDice Setup Test Script        ║${NC}"
echo -e "${BLUE}║                                           ║${NC}"
echo -e "${BLUE}╚═══════════════════════════════════════════╝${NC}"
echo ""

# Test Prerequisites
print_header "Testing Prerequisites"
echo ""

run_test "Docker installed" "command -v docker"
run_test "Go installed" "command -v go"
run_test "Node.js installed" "command -v node"
run_test "npm installed" "command -v npm"

echo ""

# Test Environment
print_header "Testing Environment Configuration"
echo ""

if [ -f .env ]; then
    print_success ".env file exists"

    # Check required variables
    if grep -q "DATABASE_URL=" .env; then
        print_success "DATABASE_URL is set"
    else
        print_error "DATABASE_URL not found in .env"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi

    if grep -q "JWT_SECRET=" .env; then
        if grep -q "JWT_SECRET=your-super-secret" .env; then
            print_warning "JWT_SECRET appears to be default value"
        else
            print_success "JWT_SECRET is configured"
        fi
    else
        print_error "JWT_SECRET not found in .env"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi

    if grep -q "OPENAI_API_KEY=" .env; then
        if grep -q "OPENAI_API_KEY=sk-your-openai-api-key" .env; then
            print_error "OPENAI_API_KEY is still the default placeholder"
            FAILED_TESTS=$((FAILED_TESTS + 1))
        else
            print_success "OPENAI_API_KEY is configured"
        fi
    else
        print_error "OPENAI_API_KEY not found in .env"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi

    TOTAL_TESTS=$((TOTAL_TESTS + 3))
else
    print_error ".env file not found"
    print_info "Run: cp .env.example .env"
    FAILED_TESTS=$((FAILED_TESTS + 1))
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
fi

echo ""

# Test Database
print_header "Testing Database Connection"
echo ""

if docker ps --format '{{.Names}}' | grep -q '^dishdice-postgres$'; then
    print_success "PostgreSQL container is running"

    # Try to connect
    if docker exec dishdice-postgres pg_isready -U postgres > /dev/null 2>&1; then
        print_success "PostgreSQL is accepting connections"
    else
        print_warning "PostgreSQL container exists but not ready yet"
    fi
else
    print_warning "PostgreSQL container is not running"
    print_info "Run: ./start.sh to start it"
fi

echo ""

# Test Backend
print_header "Testing Backend"
echo ""

if [ -f backend/go.mod ]; then
    print_success "Go module file exists"
else
    print_error "Go module file not found"
fi

if [ -f backend/go.sum ]; then
    print_success "Go dependencies downloaded"
else
    print_warning "Go dependencies not downloaded yet"
    print_info "Run: cd backend && go mod download"
fi

if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    print_success "Backend server is running"

    # Test health endpoint
    response=$(curl -s http://localhost:8080/health)
    if [ "$response" = "OK" ]; then
        print_success "Health endpoint returns OK"
    else
        print_warning "Health endpoint returned: $response"
    fi
else
    print_warning "Backend server is not running"
    print_info "Run: ./start.sh to start it"
fi

echo ""

# Test Frontend
print_header "Testing Frontend"
echo ""

if [ -f frontend/package.json ]; then
    print_success "Package.json exists"
else
    print_error "Package.json not found"
fi

if [ -d frontend/node_modules ]; then
    print_success "npm dependencies installed"
else
    print_warning "npm dependencies not installed yet"
    print_info "Run: cd frontend && npm install"
fi

if curl -s http://localhost:5173 > /dev/null 2>&1; then
    print_success "Frontend server is running"
else
    print_warning "Frontend server is not running"
    print_info "Run: ./start.sh to start it"
fi

echo ""

# Test API Connectivity
print_header "Testing API Connectivity"
echo ""

if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    # Test CORS
    cors_response=$(curl -s -I -X OPTIONS http://localhost:8080/api/auth/register \
        -H "Origin: http://localhost:5173" \
        -H "Access-Control-Request-Method: POST" 2>&1 | grep -i "access-control-allow-origin" || echo "")

    if [ -n "$cors_response" ]; then
        print_success "CORS is configured"
    else
        print_warning "CORS headers not detected (server may not be running)"
    fi
fi

echo ""

# Summary
print_header "Test Summary"
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}╔═══════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║          All Tests Passed! 🎉             ║${NC}"
    echo -e "${GREEN}╚═══════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${GREEN}Your DishDice setup is ready to use!${NC}"
    echo ""
    echo "Next steps:"
    echo "  1. If services aren't running: ./start.sh"
    echo "  2. Open browser: http://localhost:5173"
    echo "  3. Register an account and start planning meals!"
else
    echo -e "${YELLOW}╔═══════════════════════════════════════════╗${NC}"
    echo -e "${YELLOW}║          Setup Needs Attention            ║${NC}"
    echo -e "${YELLOW}╚═══════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "Passed: ${GREEN}${PASSED_TESTS}${NC}"
    echo -e "Failed: ${RED}${FAILED_TESTS}${NC}"
    echo ""

    if grep -q "OPENAI_API_KEY=sk-your-openai-api-key" .env 2>/dev/null; then
        echo -e "${YELLOW}⚠ IMPORTANT: Update your OpenAI API key in .env${NC}"
    fi

    echo ""
    echo "To fix issues:"
    echo "  1. Review the test results above"
    echo "  2. Follow the suggestions (marked with ℹ)"
    echo "  3. Run ./start.sh to start services"
    echo "  4. Run this test again: ./test-setup.sh"
fi

echo ""
