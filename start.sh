#!/bin/bash

# DishDice Local Development Startup Script
# This script starts all required services for local development

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}================================${NC}"
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

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
check_prerequisites() {
    print_header "Checking Prerequisites"

    local missing_deps=0

    # Check Docker
    if command_exists docker; then
        print_success "Docker is installed"
    else
        print_error "Docker is not installed"
        echo "  Install from: https://docs.docker.com/get-docker/"
        missing_deps=1
    fi

    # Check Go
    if command_exists go; then
        local go_version=$(go version | awk '{print $3}')
        print_success "Go is installed ($go_version)"
    else
        print_error "Go is not installed"
        echo "  Install from: https://golang.org/dl/"
        missing_deps=1
    fi

    # Check Node.js
    if command_exists node; then
        local node_version=$(node --version)
        print_success "Node.js is installed ($node_version)"
    else
        print_error "Node.js is not installed"
        echo "  Install from: https://nodejs.org/"
        missing_deps=1
    fi

    # Check npm
    if command_exists npm; then
        local npm_version=$(npm --version)
        print_success "npm is installed (v$npm_version)"
    else
        print_error "npm is not installed"
        missing_deps=1
    fi

    echo ""

    if [ $missing_deps -eq 1 ]; then
        print_error "Missing required dependencies. Please install them and try again."
        exit 1
    fi
}

# Check and create .env file
check_env_file() {
    print_header "Checking Environment Configuration"

    if [ ! -f .env ]; then
        print_warning ".env file not found"

        if [ -f .env.example ]; then
            print_info "Creating .env from .env.example..."
            cp .env.example .env
            print_success ".env file created"
            echo ""
            print_warning "IMPORTANT: Edit .env and add your OPENAI_API_KEY!"
            echo "  Open .env in your editor and replace:"
            echo "  OPENAI_API_KEY=sk-your-openai-api-key"
            echo ""
            read -p "Press Enter after you've added your OpenAI API key..."
        else
            print_error ".env.example not found. Cannot create .env"
            exit 1
        fi
    else
        print_success ".env file exists"

        # Check if OpenAI API key is set
        if grep -q "OPENAI_API_KEY=sk-your-openai-api-key" .env 2>/dev/null; then
            print_warning "OpenAI API key appears to be the default placeholder"
            echo "  Please edit .env and add your real API key"
            read -p "Press Enter after you've updated your OpenAI API key..."
        else
            print_success "OpenAI API key appears to be configured"
        fi
    fi
    echo ""
}

# Start PostgreSQL
start_postgres() {
    print_header "Starting PostgreSQL Database"

    # Check if container already exists
    if docker ps -a --format '{{.Names}}' | grep -q '^dishdice-postgres$'; then
        # Check if it's running
        if docker ps --format '{{.Names}}' | grep -q '^dishdice-postgres$'; then
            print_success "PostgreSQL container is already running"
        else
            print_info "Starting existing PostgreSQL container..."
            docker start dishdice-postgres
            print_success "PostgreSQL container started"
        fi
    else
        print_info "Creating and starting PostgreSQL container..."
        docker run --name dishdice-postgres \
            -e POSTGRES_PASSWORD=postgres \
            -e POSTGRES_DB=dishdice \
            -p 5432:5432 \
            -d postgres:16

        print_success "PostgreSQL container created and started"
        print_info "Waiting 3 seconds for PostgreSQL to initialize..."
        sleep 3
    fi
    echo ""
}

# Install backend dependencies
setup_backend() {
    print_header "Setting Up Backend"

    cd backend

    if [ ! -f go.sum ]; then
        print_info "Downloading Go dependencies..."
        go mod download
        print_success "Go dependencies downloaded"
    else
        print_success "Go dependencies already downloaded"
    fi

    cd ..
    echo ""
}

# Install frontend dependencies
setup_frontend() {
    print_header "Setting Up Frontend"

    cd frontend

    if [ ! -d node_modules ]; then
        print_info "Installing npm dependencies (this may take a minute)..."
        npm install --loglevel=error
        print_success "npm dependencies installed"
    else
        print_success "npm dependencies already installed"
    fi

    cd ..
    echo ""
}

# Create tmux session or run in separate terminals
start_services() {
    print_header "Starting Services"

    # Check if tmux is available
    if command_exists tmux; then
        print_info "Starting services in tmux session..."

        # Kill existing session if it exists
        tmux kill-session -t dishdice 2>/dev/null || true

        # Create new session
        tmux new-session -d -s dishdice -n backend

        # Start backend in first window (CGO_ENABLED=0 for macOS compatibility)
        tmux send-keys -t dishdice:backend "cd $(pwd)/backend && echo '🚀 Starting Backend Server...' && CGO_ENABLED=0 go run cmd/api/main.go" C-m

        # Create window for frontend
        tmux new-window -t dishdice -n frontend
        tmux send-keys -t dishdice:frontend "cd $(pwd)/frontend && export PATH=/opt/homebrew/bin:\$PATH && echo '🚀 Starting Frontend Server...' && npm run dev" C-m

        # Create window for logs
        tmux new-window -t dishdice -n logs
        tmux send-keys -t dishdice:logs "echo '📋 DishDice Logs'" C-m
        tmux send-keys -t dishdice:logs "echo '  Backend: http://localhost:8080'" C-m
        tmux send-keys -t dishdice:logs "echo '  Frontend: http://localhost:5173'" C-m
        tmux send-keys -t dishdice:logs "echo ''" C-m
        tmux send-keys -t dishdice:logs "echo 'Switch windows:'" C-m
        tmux send-keys -t dishdice:logs "echo '  Ctrl+b then 1: Backend logs'" C-m
        tmux send-keys -t dishdice:logs "echo '  Ctrl+b then 2: Frontend logs'" C-m
        tmux send-keys -t dishdice:logs "echo '  Ctrl+b then 3: This help'" C-m
        tmux send-keys -t dishdice:logs "echo ''" C-m
        tmux send-keys -t dishdice:logs "echo 'To exit: Ctrl+b then type :kill-session'" C-m

        # Select backend window
        tmux select-window -t dishdice:backend

        print_success "Services started in tmux session"
        echo ""
        print_info "Waiting 5 seconds for services to initialize..."
        sleep 5
        echo ""
        print_success "DishDice is ready!"
        echo ""
        echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo -e "${GREEN}  🎲 DishDice is running!${NC}"
        echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo ""
        echo -e "  ${BLUE}Frontend:${NC} http://localhost:5173"
        echo -e "  ${BLUE}Backend:${NC}  http://localhost:8080"
        echo -e "  ${BLUE}Health:${NC}   http://localhost:8080/health"
        echo ""
        echo -e "${YELLOW}To view logs, attach to tmux session:${NC}"
        echo -e "  tmux attach -t dishdice"
        echo ""
        echo -e "${YELLOW}To stop all services:${NC}"
        echo -e "  ./stop.sh"
        echo ""
        echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo ""

    else
        print_warning "tmux not found, starting services in background..."

        # Start backend (CGO_ENABLED=0 for macOS compatibility)
        cd backend
        print_info "Starting backend on http://localhost:8080..."
        CGO_ENABLED=0 go run cmd/api/main.go > ../backend.log 2>&1 &
        BACKEND_PID=$!
        echo $BACKEND_PID > ../backend.pid
        cd ..

        # Start frontend
        cd frontend
        print_info "Starting frontend on http://localhost:5173..."
        npm run dev > ../frontend.log 2>&1 &
        FRONTEND_PID=$!
        echo $FRONTEND_PID > ../frontend.pid
        cd ..

        print_success "Services started in background"
        echo ""
        print_info "Waiting 5 seconds for services to initialize..."
        sleep 5
        echo ""
        print_success "DishDice is ready!"
        echo ""
        echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo -e "${GREEN}  🎲 DishDice is running!${NC}"
        echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo ""
        echo -e "  ${BLUE}Frontend:${NC} http://localhost:5173"
        echo -e "  ${BLUE}Backend:${NC}  http://localhost:8080"
        echo -e "  ${BLUE}Health:${NC}   http://localhost:8080/health"
        echo ""
        echo -e "${YELLOW}View logs:${NC}"
        echo -e "  tail -f backend.log"
        echo -e "  tail -f frontend.log"
        echo ""
        echo -e "${YELLOW}To stop all services:${NC}"
        echo -e "  ./stop.sh"
        echo ""
        echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo ""
    fi
}

# Main execution
main() {
    clear
    echo ""
    echo -e "${BLUE}╔═══════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║                                           ║${NC}"
    echo -e "${BLUE}║       🎲 DishDice Startup Script         ║${NC}"
    echo -e "${BLUE}║                                           ║${NC}"
    echo -e "${BLUE}╚═══════════════════════════════════════════╝${NC}"
    echo ""

    check_prerequisites
    check_env_file
    start_postgres
    setup_backend
    setup_frontend
    start_services
}

# Run main function
main
