#!/bin/bash

# DishDice Deployment Script for Fly.io
# This script will deploy DishDice to Fly.io

set -e

echo "🚀 Deploying DishDice to Fly.io..."

# Check if flyctl is installed
if ! command -v flyctl &> /dev/null; then
    echo "❌ flyctl is not installed. Please install it first:"
    echo "   brew install flyctl"
    exit 1
fi

# Check if logged in
if ! flyctl auth whoami &> /dev/null; then
    echo "❌ Not logged in to Fly.io. Please run: flyctl auth login"
    exit 1
fi

echo ""
echo "Step 1: Creating Postgres database..."
echo "---------------------------------------"
if flyctl postgres list | grep -q "dishdice-db"; then
    echo "✅ Database 'dishdice-db' already exists"
else
    echo "Creating new Postgres database..."
    flyctl postgres create --name dishdice-db --region iad --initial-cluster-size 1 --vm-size shared-cpu-1x --volume-size 1
fi

echo ""
echo "Step 2: Creating app and attaching database..."
echo "-----------------------------------------------"
if flyctl apps list | grep -q "dishdice"; then
    echo "✅ App 'dishdice' already exists"
else
    echo "Creating new app..."
    flyctl launch --no-deploy --copy-config
fi

echo ""
echo "Attaching database to app..."
flyctl postgres attach dishdice-db --app dishdice || echo "Database may already be attached"

echo ""
echo "Step 3: Setting secrets..."
echo "--------------------------"
echo "Please enter your OpenAI API key:"
read -s OPENAI_API_KEY

echo ""
echo "Setting secrets..."
flyctl secrets set \
  JWT_SECRET=$(openssl rand -base64 32) \
  OPENAI_API_KEY="$OPENAI_API_KEY" \
  --app dishdice

echo ""
echo "Step 4: Deploying application..."
echo "---------------------------------"
flyctl deploy

echo ""
echo "✅ Deployment complete!"
echo ""
echo "Your app should be available at: https://dishdice.fly.dev"
echo ""
echo "To view logs: flyctl logs --app dishdice"
echo "To check status: flyctl status --app dishdice"
echo "To open the app: flyctl open --app dishdice"
