# DishDice 🎲

AI-powered weekly meal planning application that helps you decide "what should we eat this week?"

## Features

- **AI-Generated Meal Plans**: Create personalized 7-day meal plans using OpenAI GPT-4o-mini
- **Meal Regeneration**: Don't like a meal? Get 3 new AI-generated options instantly
- **Smart Shopping Lists**: Automatically generate shopping lists from meal plans
- **User Preferences**: Set dietary restrictions, favorite cuisines, and food preferences
- **Meal History Tracking**: AI avoids repeating recent meals for variety
- **Resource Utilization**: Tell the AI what ingredients you have, and it will incorporate them
- **Admin Approval System**: First user becomes admin, subsequent users require approval
- **Pre-Approved Registration Links**: Admins can generate single-use registration links that auto-approve users

## Tech Stack

### Backend
- **Go** with Chi router
- **PostgreSQL** database
- **OpenAI GPT-4o-mini** for AI meal generation
- **JWT** authentication

### Frontend
- **React** with TypeScript
- **Vite** for fast development
- **Tailwind CSS** for styling
- **React Router** for navigation
- **Axios** for API calls

## Prerequisites

- Go 1.22+
- Node.js 18+
- PostgreSQL 14+
- OpenAI API key

## Local Development Setup

### Quick Start (Recommended)

Use the automated startup script:

```bash
./start.sh
```

This will:
- ✅ Check prerequisites (Docker, Go, Node.js)
- ✅ Create .env from template
- ✅ Start PostgreSQL in Docker
- ✅ Install all dependencies
- ✅ Start backend and frontend servers

**First time setup:** The script will prompt you to add your OpenAI API key to `.env`

**Stop services:**
```bash
./stop.sh
```

See [SCRIPTS.md](SCRIPTS.md) for detailed script documentation.

### Manual Setup

If you prefer manual setup:

#### 1. Clone the repository

```bash
git clone <your-repo-url>
cd dishdice
```

#### 2. Set up environment variables

Create a `.env` file in the root directory:

```bash
cp .env.example .env
```

Edit `.env` with your values:
```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/dishdice?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-change-this
OPENAI_API_KEY=sk-your-openai-api-key
PORT=8080
ALLOWED_ORIGINS=http://localhost:5173
```

#### 3. Set up the database

```bash
# Start PostgreSQL (if using Docker)
docker run --name dishdice-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=dishdice -p 5432:5432 -d postgres:16

# Migrations will run automatically when starting the backend
```

#### 4. Start the backend

```bash
cd backend
go mod download
CGO_ENABLED=0 go run cmd/api/main.go  # CGO_ENABLED=0 for macOS compatibility
```

The backend will run on `http://localhost:8080`

#### 5. Start the frontend

```bash
cd frontend
npm install
npm run dev
```

The frontend will run on `http://localhost:5173`

## Usage

1. **Register/Login**: Create an account or log in
2. **Set Preferences**: Go to Preferences and add your dietary restrictions, favorite cuisines, etc.
3. **Create Meal Plan**:
   - Click "New Proposal"
   - Select the week start date
   - Optionally add week-specific preferences (e.g., "I want pasta on Friday")
   - Optionally list available ingredients
   - Click "Generate Meal Plan"
4. **Manage Meals**:
   - View your 7-day meal plan
   - Click "Regenerate" on any day to get 3 new options
   - Select your preferred option
5. **Shopping List**:
   - Click "Add to List" on any meal to add its ingredients
   - Manually add items
   - Check off items as you shop
   - Clear completed items

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login
- `GET /api/auth/me` - Get current user

### User Preferences
- `GET /api/user/preferences` - Get preferences
- `PUT /api/user/preferences` - Update preferences

### Proposals
- `GET /api/proposals` - List proposals
- `POST /api/proposals` - Create proposal
- `GET /api/proposals/:id` - Get proposal details
- `DELETE /api/proposals/:id` - Delete proposal

### Meals
- `POST /api/meals/:id/regenerate` - Generate 3 new options
- `PUT /api/meals/:id/select` - Select an option
- `POST /api/meals/:id/save-to-shopping` - Add to shopping list

### Shopping List
- `GET /api/shopping-list` - Get items
- `POST /api/shopping-list` - Add item
- `PUT /api/shopping-list/:id/toggle` - Toggle checked
- `DELETE /api/shopping-list/checked` - Delete checked items
- `DELETE /api/shopping-list/:id` - Delete item

## Deployment

### Fly.io Deployment

1. **Install Fly CLI**:
```bash
curl -L https://fly.io/install.sh | sh
```

2. **Create Fly.io apps**:
```bash
fly launch --no-deploy
```

3. **Create PostgreSQL**:
```bash
fly postgres create --name dishdice-db
fly postgres attach dishdice-db
```

4. **Set secrets**:
```bash
fly secrets set JWT_SECRET=your-secret-key
fly secrets set OPENAI_API_KEY=sk-your-key
fly secrets set ALLOWED_ORIGINS=https://your-app.fly.dev
```

Note: `FRONTEND_URL` is already set in `fly.toml` to `https://dishdice.fly.dev`

5. **Deploy**:
```bash
fly deploy
```

## Cost Estimate

- **Fly.io**: $0 (free tier with auto-stop/start)
- **Fly Postgres**: $0 (256MB RAM, 1GB storage)
- **OpenAI**: ~$1-3/month (GPT-4o-mini at $0.15/1M tokens)
- **Total**: $1-3/month

## Project Structure

```
dishdice/
├── backend/
│   ├── cmd/api/main.go           # Application entry point
│   ├── internal/
│   │   ├── ai/                   # OpenAI integration
│   │   ├── config/               # Configuration
│   │   ├── database/             # Database & migrations
│   │   ├── handlers/             # HTTP handlers
│   │   ├── middleware/           # Auth & logging middleware
│   │   ├── models/               # Data models
│   │   ├── repository/           # Database queries
│   │   └── services/             # Business logic
│   └── migrations/               # SQL migrations
├── frontend/
│   └── src/
│       ├── components/           # React components
│       ├── context/              # Auth context
│       ├── pages/                # Page components
│       ├── services/             # API services
│       └── types/                # TypeScript types
└── README.md
```

## License

MIT
