#!/bin/bash

# DishDice Stop Script
# This script stops all DishDice services

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

main() {
    clear
    echo ""
    print_header "Stopping DishDice Services"
    echo ""

    # Stop tmux session if it exists
    if command_exists tmux && tmux has-session -t dishdice 2>/dev/null; then
        print_info "Stopping tmux session..."
        tmux kill-session -t dishdice
        print_success "tmux session stopped"
    fi

    # Stop background processes if PIDs exist
    if [ -f backend.pid ]; then
        print_info "Stopping backend process..."
        kill $(cat backend.pid) 2>/dev/null || true
        rm backend.pid
        print_success "Backend stopped"
    fi

    if [ -f frontend.pid ]; then
        print_info "Stopping frontend process..."
        kill $(cat frontend.pid) 2>/dev/null || true
        rm frontend.pid
        print_success "Frontend stopped"
    fi

    # Ask about PostgreSQL
    echo ""
    read -p "Stop PostgreSQL container? (y/N): " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        if command_exists docker && docker ps --format '{{.Names}}' | grep -q '^dishdice-postgres$'; then
            print_info "Stopping PostgreSQL container..."
            docker stop dishdice-postgres
            print_success "PostgreSQL container stopped"
        else
            print_info "PostgreSQL container is not running"
        fi
    else
        print_info "PostgreSQL container left running"
    fi

    echo ""
    print_success "DishDice services stopped"
    echo ""
}

main
