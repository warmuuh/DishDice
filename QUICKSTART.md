# DishDice Quick Start Guide

Get DishDice running locally in 5 minutes!

## Prerequisites

- Docker (for PostgreSQL)
- Go 1.22+
- Node.js 18+
- OpenAI API key (get one at https://platform.openai.com/)

## Setup Steps

### 1. Start PostgreSQL

```bash
docker run --name dishdice-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=dishdice \
  -p 5432:5432 \
  -d postgres:16
```

### 2. Configure Environment

Create `.env` file in the root directory:

```bash
DATABASE_URL=postgres://postgres:postgres@localhost:5432/dishdice?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-change-this-to-random-string
OPENAI_API_KEY=sk-your-openai-api-key-here
PORT=8080
ALLOWED_ORIGINS=http://localhost:5173
```

**Important**: Replace `OPENAI_API_KEY` with your actual OpenAI API key!

### 3. Start Backend

```bash
cd backend
go mod download
CGO_ENABLED=0 go run cmd/api/main.go  # CGO_ENABLED=0 fixes macOS build issues
```

You should see:
```
Connected to database successfully
Migrations completed successfully
Server starting on :8080
```

### 4. Start Frontend (New Terminal)

```bash
cd frontend
npm install
npm run dev
```

You should see:
```
VITE v8.0.0  ready in XXX ms

➜  Local:   http://localhost:5173/
```

### 5. Use the App

1. Open http://localhost:5173 in your browser
2. Click "Sign up" and create an account
3. Go to "Preferences" and add your food preferences:
   ```
   I love Italian food, prefer vegetarian meals.
   No shellfish allergies. Keep it family-friendly and not too spicy.
   ```
4. Click "New Proposal" from the Dashboard
5. Select next Monday's date
6. Optionally add:
   - Week preferences: "I want lasagna on Friday"
   - Available ingredients: "I have gouda cheese, carrots, and chicken"
7. Click "Generate Meal Plan" and wait 10-20 seconds
8. Enjoy your AI-generated 7-day meal plan!

## Features to Try

### Regenerate a Meal
- Click "Regenerate" on any day card
- Get 3 new AI-generated options
- Select your favorite

### Shopping List
- Click "Add to List" on any meal
- Items are automatically added to your shopping list
- Navigate to "Shopping List" in the header
- Check off items as you shop
- Click "Clear Checked" when done

### Multiple Proposals
- Create multiple proposals for different weeks
- View all proposals on the Dashboard
- Click any proposal to view its 7-day plan

## Troubleshooting

### Backend won't start
- Check PostgreSQL is running: `docker ps`
- Verify DATABASE_URL in .env
- Check you have an OpenAI API key set

### Frontend won't build
- Delete node_modules and reinstall: `rm -rf node_modules && npm install`
- Clear Vite cache: `rm -rf .vite`

### AI not generating meals
- Verify your OpenAI API key is valid
- Check you have API credits
- Look at backend logs for error messages

### Database errors
- Reset database: `docker rm -f dishdice-postgres` then restart from step 1
- Migrations run automatically on backend start

## API Testing with curl

### Register
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### Create Proposal (with token)
```bash
curl -X POST http://localhost:8080/api/proposals \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "week_start_date": "2026-03-17",
    "week_preferences": "Light meals",
    "current_resources": "Chicken and vegetables"
  }'
```

## Development Tips

- Backend auto-reloads with: `go run cmd/api/main.go`
- Frontend has hot reload built-in with Vite
- Check backend logs for AI prompt/response debugging
- Use browser DevTools Network tab to debug API calls

## Next Steps

- Read the full [README.md](README.md) for detailed documentation
- See [IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md) for technical details
- Deploy to Fly.io (instructions in README.md)

## Need Help?

- Backend logs: Check the terminal running `go run cmd/api/main.go`
- Frontend errors: Check browser console (F12)
- API errors: Check Network tab in browser DevTools

Happy meal planning! 🎲
