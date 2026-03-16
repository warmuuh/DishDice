# DishDice Deployment Guide for Fly.io

## Prerequisites

1. **Add Payment Information**: Go to https://fly.io/dashboard/peter-866/billing and add a credit card
2. **Verify Login**: Run `flyctl auth whoami` to confirm you're logged in
3. **OpenAI API Key**: Have your OpenAI API key ready

## Option 1: Automated Deployment (Recommended)

Simply run the deployment script:

```bash
./deploy.sh
```

The script will:
1. Create a Postgres database
2. Create and configure the app
3. Set required secrets
4. Deploy the application

## Option 2: Manual Deployment

### Step 1: Create Postgres Database

```bash
flyctl postgres create --name dishdice-db --region iad --initial-cluster-size 1 --vm-size shared-cpu-1x --volume-size 1
```

### Step 2: Create App

```bash
flyctl launch --no-deploy --copy-config
```

### Step 3: Attach Database

```bash
flyctl postgres attach dishdice-db --app dishdice
```

This will automatically set the `DATABASE_URL` secret.

### Step 4: Set Secrets

```bash
# Generate a random JWT secret
JWT_SECRET=$(openssl rand -base64 32)

# Set secrets
flyctl secrets set \
  JWT_SECRET="$JWT_SECRET" \
  OPENAI_API_KEY="your-openai-api-key-here" \
  --app dishdice
```

### Step 5: Deploy

```bash
flyctl deploy
```

## Post-Deployment

### Check Status

```bash
flyctl status --app dishdice
```

### View Logs

```bash
flyctl logs --app dishdice
```

### Open the App

```bash
flyctl open --app dishdice
```

Your app will be available at: **https://dishdice.fly.dev**

## Environment Variables

The following environment variables are automatically configured:

- `DATABASE_URL` - Set by Postgres attachment
- `JWT_SECRET` - Set via secrets
- `OPENAI_API_KEY` - Set via secrets
- `PORT` - Set to 8080 in fly.toml
- `ALLOWED_ORIGINS` - Will be set to https://dishdice.fly.dev

## Database Migrations

Migrations run automatically on startup. The backend will:
1. Check which migrations have been applied
2. Run any pending migrations
3. Start the server

## Cost Optimization

The app is configured for cost optimization:

- **Auto-stop**: Machines stop when idle (no traffic)
- **Auto-start**: Machines start when traffic arrives
- **Min machines**: 0 (only runs when needed)
- **VM size**: 256MB RAM (smallest size)
- **Postgres**: 1GB volume (smallest size)

Expected monthly cost: **$0-5** (mostly for Postgres storage)

## Troubleshooting

### App Won't Start

Check logs:
```bash
flyctl logs --app dishdice
```

Common issues:
- Missing secrets (DATABASE_URL, JWT_SECRET, OPENAI_API_KEY)
- Database connection issues
- Migration failures

### Database Connection Issues

Check database status:
```bash
flyctl postgres db list --app dishdice-db
```

Get connection string:
```bash
flyctl postgres db show --app dishdice-db
```

### Reset Everything

If you need to start over:

```bash
# Delete app
flyctl apps destroy dishdice

# Delete database
flyctl postgres destroy dishdice-db

# Start fresh
./deploy.sh
```

## Monitoring

### Health Check

The app has a health endpoint at `/health` that Fly.io checks every 30 seconds.

### Scaling

If you need more resources:

```bash
# Scale memory
flyctl scale memory 512 --app dishdice

# Add more machines
flyctl scale count 2 --app dishdice
```

## Updating the App

After making code changes:

```bash
# Commit changes
git add .
git commit -m "Your changes"

# Deploy
flyctl deploy
```

## Support

- Fly.io Docs: https://fly.io/docs/
- DishDice GitHub: Create an issue
- OpenAI API: https://platform.openai.com/docs
