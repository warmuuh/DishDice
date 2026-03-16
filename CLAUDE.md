# DishDice - Claude Code Project Guide

This file helps Claude Code understand the DishDice project structure, conventions, and important context for future sessions.

## Project Overview

**DishDice** is an AI-powered weekly meal planning application that helps users decide "what should we eat this week?" by generating personalized 7-day meal plans using OpenAI's GPT-4o-mini.

### Core Value Proposition
- Users set dietary preferences once
- AI generates complete 7-day meal plans with recipes and shopping lists
- Users can regenerate individual days to get 3 new options
- Shopping lists are automatically generated from meal plans
- AI tracks meal history to ensure variety (no repetition)

### Tech Stack
- **Backend**: Go 1.22+ with Chi router
- **Database**: PostgreSQL 14+ with custom migration system
- **AI**: OpenAI GPT-4o-mini API
- **Frontend**: React 18+ with TypeScript, Vite, Tailwind CSS
- **Auth**: JWT tokens with bcrypt password hashing
- **Deployment**: Fly.io with auto-stop/start for cost optimization

## Project Status

✅ **FULLY IMPLEMENTED** - All features from the original plan are complete and working:
- Authentication and user management
- AI meal plan generation (7 days)
- Meal regeneration with 3 options
- Shopping list management
- User preferences
- Responsive UI with Tailwind CSS
- Production-ready Dockerfile and Fly.io configuration

## Architecture

### Backend Structure (Go)

```
backend/
├── cmd/api/main.go              # Application entry point
│   └── Initializes all services, runs migrations, starts server
├── internal/
│   ├── ai/                      # OpenAI Integration (KEY FEATURE)
│   │   ├── client.go            # OpenAI API calls with structured output
│   │   ├── prompts.go           # Prompt engineering for meal generation
│   │   └── types.go             # AI request/response types
│   ├── config/config.go         # Environment variable management
│   ├── database/
│   │   ├── db.go                # PostgreSQL connection pooling
│   │   └── migrations.go        # Custom migration runner
│   ├── handlers/                # HTTP request handlers
│   │   ├── auth_handler.go      # Register, login, get-me
│   │   ├── user_handler.go      # Get/update preferences
│   │   ├── proposal_handler.go  # CRUD for proposals
│   │   ├── meal_handler.go      # Regenerate, select, add-to-shopping
│   │   └── shopping_handler.go  # Shopping list CRUD
│   ├── middleware/
│   │   ├── auth.go              # JWT validation
│   │   └── logging.go           # Request logging
│   ├── models/                  # Domain models and DTOs
│   │   ├── user.go
│   │   ├── proposal.go
│   │   └── shopping.go
│   ├── repository/              # Database access layer
│   │   ├── user_repository.go
│   │   ├── proposal_repository.go
│   │   └── shopping_repository.go
│   └── services/                # Business logic layer
│       ├── auth_service.go      # Authentication + JWT
│       ├── proposal_service.go  # Orchestrates AI + DB for proposals
│       ├── meal_service.go      # Meal regeneration logic
│       └── shopping_service.go  # Shopping list operations
└── migrations/
    ├── 001_init.up.sql          # Initial schema
    └── 001_init.down.sql        # Rollback schema
```

### Frontend Structure (React + TypeScript)

```
frontend/src/
├── components/
│   ├── DayCard.tsx              # Individual day meal display
│   ├── Header.tsx               # Navigation with gradient
│   ├── LoadingSpinner.tsx       # Loading state
│   ├── ProtectedRoute.tsx       # Auth guard wrapper
│   ├── RegenerateModal.tsx      # Modal showing 3 AI options
│   └── ShoppingListItem.tsx     # Checkable shopping item
├── context/
│   └── AuthContext.tsx          # Global auth state (user, token, login, logout)
├── pages/
│   ├── Dashboard.tsx            # List of all proposals
│   ├── Login.tsx                # Login form
│   ├── NewProposal.tsx          # Create new meal plan
│   ├── Preferences.tsx          # User food preferences
│   ├── ProposalDetail.tsx       # 7-day meal plan view (KEY PAGE)
│   ├── Register.tsx             # Registration form
│   └── ShoppingList.tsx         # Shopping list management
├── services/                    # API client layer
│   ├── api.ts                   # Axios instance with auth interceptor
│   ├── authService.ts
│   ├── mealService.ts
│   ├── proposalService.ts
│   ├── shoppingService.ts
│   └── userService.ts
├── types/index.ts               # TypeScript type definitions
├── App.tsx                      # React Router setup
├── index.css                    # Tailwind directives + Google Fonts
└── main.tsx                     # React entry point
```

## Database Schema

### Key Tables
1. **users** - User accounts with email, password_hash, preferences
2. **weekly_proposals** - Weekly meal plans with start date and context
3. **daily_meals** - Individual meals (7 per proposal) with recipes
4. **meal_shopping_items** - Ingredients for each meal
5. **shopping_list** - User's shopping list (manual + from meals)
6. **meal_generation_history** - Tracks recent meals to avoid repetition

### Important Relationships
- Users → Weekly Proposals (1:many)
- Weekly Proposals → Daily Meals (1:7)
- Daily Meals → Meal Shopping Items (1:many)
- Daily Meals → Shopping List (optional, when user adds meal to list)

### Indexes
- `idx_proposals_user_date` - Fast proposal lookup by user and date
- `idx_history_user` - Quick access to recent meal history

## Critical Workflows

### 1. Create Weekly Proposal (Main Feature)
**Path**: `POST /api/proposals`

**Flow**:
1. User submits: week_start_date, week_preferences (optional), current_resources (optional)
2. Backend fetches user's general preferences
3. Backend queries last 20 meals from generation history
4. Backend calls OpenAI with structured prompt including:
   - User preferences
   - Week preferences
   - Available ingredients
   - Recent meals to avoid
5. OpenAI returns JSON with 7 days of meals (name, recipe, shopping items)
6. Backend begins database transaction:
   - Creates weekly_proposal record
   - Creates 7 daily_meal records
   - Creates meal_shopping_items for each meal
   - Adds meal names to generation history
   - Commits transaction
7. Returns complete proposal to frontend

**Key File**: `backend/internal/services/proposal_service.go:CreateWeeklyProposal()`

### 2. Regenerate Single Meal
**Path**: `POST /api/meals/:id/regenerate`

**Flow**:
1. User clicks "Regenerate" on a day card
2. Backend fetches the meal and its proposal
3. Backend gathers context:
   - User preferences
   - Week preferences
   - Other 6 meals in the week
   - Recent meal history
4. Backend calls OpenAI asking for 3 diverse alternatives
5. OpenAI returns 3 complete meal options
6. Frontend displays modal with 3 options
7. User selects one, calls `PUT /api/meals/:id/select`
8. Backend updates meal in transaction, adds to history

**Key Files**:
- `backend/internal/services/meal_service.go:RegenerateDayOptions()`
- `frontend/src/pages/ProposalDetail.tsx`
- `frontend/src/components/RegenerateModal.tsx`

### 3. Add Meal to Shopping List
**Path**: `POST /api/meals/:id/save-to-shopping`

**Flow**:
1. User clicks "Add to List" on a day card
2. Backend fetches meal's shopping items
3. Backend copies each item to shopping_list table with source='meal'
4. Frontend shows success toast
5. User can view items in Shopping List page

## Environment Variables

### Required Variables
```bash
DATABASE_URL=postgres://user:password@host:5432/dishdice?sslmode=disable
JWT_SECRET=random-secure-string-at-least-32-chars
OPENAI_API_KEY=sk-your-openai-api-key
PORT=8080  # Optional, defaults to 8080
ALLOWED_ORIGINS=http://localhost:5173  # For CORS
```

### Frontend Environment (optional)
```bash
VITE_API_URL=http://localhost:8080  # Defaults to this if not set
```

## Development Commands

### Backend
```bash
cd backend
go mod download              # Install dependencies
go run cmd/api/main.go       # Run dev server
go build -o dishdice cmd/api/main.go  # Build binary
```

### Frontend
```bash
cd frontend
npm install                  # Install dependencies
npm run dev                  # Run dev server (Vite)
npm run build               # Build for production
npm run preview             # Preview production build
```

### Database
```bash
# Start PostgreSQL (Docker)
docker run --name dishdice-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=dishdice \
  -p 5432:5432 -d postgres:16

# Migrations run automatically on backend startup
# No manual migration commands needed
```

## Code Conventions

### Backend (Go)
- **Package naming**: lowercase, singular (e.g., `repository`, not `repositories`)
- **Error handling**: Always return errors, wrap with context using `fmt.Errorf("context: %w", err)`
- **Transactions**: Use `sql.Tx` for multi-operation writes
- **SQL**: Use parameterized queries ($1, $2) to prevent injection
- **Handlers**: Return JSON with proper status codes
- **Models**: Use pointers for nullable fields (e.g., `*string` for optional text)

### Frontend (React)
- **Components**: Functional components with TypeScript
- **Hooks**: Use `useState`, `useEffect`, custom hooks for logic
- **API calls**: Always use try/catch with toast notifications
- **Styling**: Tailwind CSS utility classes, custom colors in tailwind.config.js
- **Types**: Import from `src/types/index.ts`
- **Loading states**: Show spinners during async operations

### API Conventions
- **Auth**: Protected routes require `Authorization: Bearer <token>` header
- **Errors**: Return descriptive error messages in response body
- **Success**: Use appropriate status codes (200, 201, 204)
- **Validation**: Check required fields before processing

## AI Prompt Engineering

### Weekly Plan Prompt Strategy
- Include user's general preferences first
- Add week-specific requests
- List available ingredients
- Show recent meal history to avoid
- Request variety (no repeated proteins/cuisines)
- Use temperature 0.8 for creativity
- Demand JSON format with exact structure

### Day Options Prompt Strategy
- Same context as weekly prompt
- Include OTHER 6 days already planned
- Request 3 DIVERSE options (different cuisines/styles)
- Maintain consistency with week preferences

**Key Insight**: The AI is very good at following constraints when given clear instructions and structured output format.

## Common Issues & Solutions

### Issue: AI Not Returning Valid JSON
**Solution**: OpenAI's `response_format: json_object` ensures valid JSON. If parsing fails, check the prompt template for typos.

### Issue: Meals Repeating Too Often
**Solution**: The `meal_generation_history` table tracks last 20 meals. Increase the limit in `GetMealHistory()` if needed.

### Issue: CORS Errors in Frontend
**Solution**: Check `ALLOWED_ORIGINS` in backend .env matches frontend URL. Both `http://localhost:5173` and production URL should be allowed.

### Issue: JWT Token Expired
**Solution**: Tokens expire after 7 days. Frontend auto-redirects to login on 401. User must log in again.

### Issue: Database Migration Fails
**Solution**: Migrations run on startup. If schema changes, drop the DB and restart, or create a new migration file (002_*.up.sql).

## Security Considerations

### Implemented
✅ Password hashing with bcrypt (cost 12)
✅ JWT token validation on all protected routes
✅ Authorization checks (user can only access their own data)
✅ Parameterized SQL queries (no injection)
✅ CORS properly configured
✅ Environment variables for secrets

### DO NOT
❌ Never commit .env files
❌ Never log sensitive data (passwords, tokens)
❌ Never trust client-side data without validation
❌ Never expose internal error details to users

## Cost Optimization

### Current Setup (Goal: <$5/month)
- **Fly.io**: $0 - Uses auto-stop/start with `min_machines_running = 0`
- **Fly Postgres**: $0 - Free tier (256MB RAM, 1GB storage)
- **OpenAI**: ~$1-3/month - GPT-4o-mini is very cheap ($0.15/1M tokens)

### Monitoring
- Check Fly.io dashboard for machine usage
- Monitor OpenAI usage at platform.openai.com
- Database size stays small (meal plans are text-only)

## Testing Strategy (Not Yet Implemented)

### Recommended Tests
1. **Unit Tests**: Services layer (auth, proposal generation logic)
2. **Integration Tests**: API endpoints with test database
3. **E2E Tests**: Complete user flows (register → create proposal → shop)

### Test Data
- Use factory functions for creating test users/proposals
- Mock OpenAI responses for consistent AI behavior
- Seed database with sample meal plans for UI testing

## Future Enhancements (v2 Ideas)

### High Priority
- [ ] Recipe image generation (DALL-E)
- [ ] Nutrition information (calories, macros)
- [ ] Email notifications (weekly reminders)
- [ ] Share proposals with family members

### Medium Priority
- [ ] Favorite/rate meals for better recommendations
- [ ] Meal prep mode (batch cooking)
- [ ] Grocery store aisle organization for shopping list
- [ ] Recipe notes and modifications

### Low Priority
- [ ] Mobile app (React Native)
- [ ] Social features (share/discover recipes)
- [ ] Meal history charts/analytics
- [ ] Multiple dietary profiles per user

## Deployment

### Fly.io Production
```bash
# One-time setup
fly launch --no-deploy
fly postgres create --name dishdice-db
fly postgres attach dishdice-db

# Set secrets
fly secrets set JWT_SECRET=<random-string>
fly secrets set OPENAI_API_KEY=<your-key>
fly secrets set ALLOWED_ORIGINS=https://your-app.fly.dev

# Deploy
fly deploy

# Monitor
fly logs
fly status
```

### Docker Local Testing
```bash
docker build -t dishdice .
docker run -p 8080:8080 \
  -e DATABASE_URL=$DATABASE_URL \
  -e JWT_SECRET=$JWT_SECRET \
  -e OPENAI_API_KEY=$OPENAI_API_KEY \
  dishdice
```

## Key Files to Know

### Backend
- `cmd/api/main.go` - Start here, wires everything together
- `internal/ai/client.go` - OpenAI integration, core differentiator
- `internal/services/proposal_service.go` - Main business logic
- `internal/database/migrations.go` - Custom migration system

### Frontend
- `src/pages/ProposalDetail.tsx` - Main UX, 7-day view
- `src/context/AuthContext.tsx` - Auth state management
- `src/services/api.ts` - API client with interceptors

### Deployment
- `Dockerfile` - Multi-stage build (frontend + backend)
- `fly.toml` - Fly.io configuration with auto-stop

## Documentation Files
- `README.md` - Full project documentation
- `QUICKSTART.md` - 5-minute setup guide
- `IMPLEMENTATION_SUMMARY.md` - Technical implementation details
- `CLAUDE.md` - This file (project guide for Claude)

## Important Notes for Claude

1. **Never generate URLs** unless they're for actual programming (APIs, webhooks, etc.). The app doesn't need external URLs beyond OpenAI API.

2. **AI is the core feature** - Any changes to prompt engineering in `internal/ai/prompts.go` should be tested carefully as it affects the entire UX.

3. **Transactions are critical** - Proposal creation must be atomic. If any step fails, rollback everything.

4. **Cost awareness** - Always consider OpenAI token usage when modifying prompts. Longer prompts = higher costs.

5. **User authorization** - Every endpoint that accesses user data MUST check that the authenticated user owns that data.

6. **Frontend state management** - Auth context is the only global state. Everything else is local to pages/components.

7. **No over-engineering** - This is a simple app. Don't add unnecessary abstractions, caching, or optimization until there's a proven need.

## Contact & Support

- This is a demonstration project for the implementation plan
- For bugs or questions about the implementation, refer to the code comments and documentation
- OpenAI API documentation: https://platform.openai.com/docs
- Fly.io documentation: https://fly.io/docs

---

**Last Updated**: March 16, 2026
**Version**: 1.0.0 (Initial Release)
**Status**: Production Ready ✅
