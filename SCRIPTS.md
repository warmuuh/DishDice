# DishDice Startup Scripts

This directory contains convenient bash scripts to manage the DishDice application locally.

## Scripts

### `test-setup.sh` - Test Your Setup

Diagnostic script that verifies everything is configured correctly:
- ✅ Checks all prerequisites
- ✅ Validates .env configuration
- ✅ Tests database connectivity
- ✅ Verifies backend/frontend status
- ✅ Tests API endpoints and CORS

**Usage:**
```bash
./test-setup.sh
```

Run this anytime to verify your setup or troubleshoot issues.

### `start.sh` - Start All Services

Automated startup script that:
- ✅ Checks prerequisites (Docker, Go, Node.js)
- ✅ Validates .env configuration
- ✅ Starts PostgreSQL in Docker
- ✅ Installs backend dependencies
- ✅ Installs frontend dependencies
- ✅ Starts backend and frontend servers

**Usage:**
```bash
./start.sh
```

**With tmux (Recommended):**
If tmux is installed, services run in a managed session with separate windows:
- Window 1: Backend logs
- Window 2: Frontend logs
- Window 3: Help/info

Access the session:
```bash
tmux attach -t dishdice
```

Navigate between windows:
- `Ctrl+b` then `1` - Backend window
- `Ctrl+b` then `2` - Frontend window
- `Ctrl+b` then `3` - Help window

Exit tmux: `Ctrl+b` then type `:kill-session`

**Without tmux:**
Services run in background, logs written to files:
- `backend.log` - Backend output
- `frontend.log` - Frontend output

### `stop.sh` - Stop All Services

Gracefully stops all DishDice services:
- ✅ Stops tmux session (if running)
- ✅ Stops background processes (if running)
- ✅ Optionally stops PostgreSQL container

**Usage:**
```bash
./stop.sh
```

You'll be prompted whether to stop PostgreSQL (useful if you want to keep data between runs).

## First Time Setup

### 1. Install Prerequisites

**macOS:**
```bash
# Install Homebrew (if not installed)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install dependencies
brew install docker go node tmux  # tmux is optional but recommended

# Start Docker Desktop
open -a Docker
```

**Linux (Ubuntu/Debian):**
```bash
# Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Go
sudo apt install golang-go

# Node.js
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt-get install -y nodejs

# tmux (optional)
sudo apt install tmux
```

**Windows (WSL2):**
```bash
# Install Docker Desktop for Windows with WSL2 backend
# Then in WSL2:
sudo apt update
sudo apt install golang-go nodejs npm tmux
```

### 2. Get OpenAI API Key

1. Go to https://platform.openai.com/
2. Sign up or log in
3. Navigate to API Keys section
4. Create a new secret key
5. Copy the key (starts with `sk-`)

### 3. Run Startup Script

```bash
./start.sh
```

On first run, the script will:
- Create `.env` from `.env.example`
- Prompt you to add your OpenAI API key
- Set up all dependencies
- Start all services

## URLs After Startup

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health

## Common Operations

### View Logs (tmux)
```bash
# Attach to session
tmux attach -t dishdice

# Switch windows
Ctrl+b then 1  # Backend
Ctrl+b then 2  # Frontend
Ctrl+b then 3  # Help
```

### View Logs (background mode)
```bash
# Real-time logs
tail -f backend.log
tail -f frontend.log

# Last 50 lines
tail -n 50 backend.log
```

### Restart Services
```bash
./stop.sh
./start.sh
```

### Reset Database
```bash
# Stop PostgreSQL
docker stop dishdice-postgres
docker rm dishdice-postgres

# Restart everything (creates fresh DB)
./start.sh
```

### Update Dependencies
```bash
# Backend
cd backend
go mod tidy
go mod download

# Frontend
cd frontend
npm install
```

## Troubleshooting

### "Port already in use"
Another service is using port 8080 or 5173:
```bash
# Find what's using the port
lsof -i :8080
lsof -i :5173

# Kill the process
kill -9 <PID>
```

### "Docker command not found"
Docker isn't installed or not running:
```bash
# macOS: Start Docker Desktop
open -a Docker

# Linux: Start Docker service
sudo systemctl start docker
```

### "PostgreSQL connection refused"
Database isn't ready yet:
```bash
# Wait a few seconds and try again
# Or restart PostgreSQL
docker restart dishdice-postgres
```

### "OpenAI API error"
Check your API key:
```bash
# Verify key is set
grep OPENAI_API_KEY .env

# Should see: OPENAI_API_KEY=sk-...
# If it says "sk-your-openai-api-key", update with real key
```

### Frontend build errors
Clear cache and reinstall:
```bash
cd frontend
rm -rf node_modules .vite dist
npm install
cd ..
./start.sh
```

### Backend build errors
Clear Go cache:
```bash
cd backend
go clean -cache
go mod tidy
cd ..
./start.sh
```

## Script Features

### Automatic Dependency Installation
- ✅ Detects if dependencies are already installed
- ✅ Skips reinstallation if present
- ✅ Fast subsequent startups

### Smart Container Management
- ✅ Reuses existing PostgreSQL container
- ✅ Doesn't recreate if already running
- ✅ Preserves database data between runs

### Environment Validation
- ✅ Checks for required tools
- ✅ Validates .env configuration
- ✅ Warns about placeholder values

### Graceful Shutdown
- ✅ Stops services cleanly
- ✅ Optional PostgreSQL shutdown
- ✅ Cleans up PID files

## Advanced Usage

### Custom PostgreSQL Port
Edit the script to use a different port:
```bash
# In start.sh, change:
-p 5432:5432
# To:
-p 5433:5432

# Then update .env:
DATABASE_URL=postgres://postgres:postgres@localhost:5433/dishdice?sslmode=disable
```

### Run Without Script
Manual startup (not recommended):
```bash
# Terminal 1: PostgreSQL
docker run --name dishdice-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=dishdice -p 5432:5432 -d postgres:16

# Terminal 2: Backend
cd backend && go run cmd/api/main.go

# Terminal 3: Frontend
cd frontend && npm run dev
```

### Debug Mode
Add debug flags for verbose output:
```bash
# Backend with verbose logs
cd backend
go run cmd/api/main.go -v

# Frontend with debug info
cd frontend
npm run dev -- --debug
```

## Environment Variables

The startup script checks these in `.env`:

**Required:**
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - Secret key for JWT tokens (min 32 chars)
- `OPENAI_API_KEY` - Your OpenAI API key

**Optional:**
- `PORT` - Backend port (default: 8080)
- `ALLOWED_ORIGINS` - CORS origins (default: http://localhost:5173)

## Tips

1. **Use tmux** - Much better experience than background processes
2. **Keep PostgreSQL running** - Faster restarts, preserves data
3. **Watch logs** - Helpful for debugging issues
4. **Check health endpoint** - Verify backend is responding: `curl http://localhost:8080/health`
5. **Use stop script** - Cleaner than Ctrl+C

## Support

If scripts don't work:
1. Check prerequisites are installed
2. Verify .env has real OpenAI API key
3. Check ports 5173 and 8080 are available
4. Look at log files for errors
5. Try manual startup to isolate the issue

For more help, see:
- [README.md](README.md) - Full documentation
- [QUICKSTART.md](QUICKSTART.md) - Manual setup guide
- [CLAUDE.md](CLAUDE.md) - Technical details
