-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    preferences TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create weekly_proposals table
CREATE TABLE IF NOT EXISTS weekly_proposals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    week_start_date DATE NOT NULL,
    week_preferences TEXT,
    current_resources TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, week_start_date)
);

-- Create daily_meals table
CREATE TABLE IF NOT EXISTS daily_meals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    proposal_id UUID NOT NULL REFERENCES weekly_proposals(id) ON DELETE CASCADE,
    day_of_week INT NOT NULL CHECK (day_of_week >= 0 AND day_of_week <= 6),
    menu_name VARCHAR(255) NOT NULL,
    recipe TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(proposal_id, day_of_week)
);

-- Create meal_shopping_items table
CREATE TABLE IF NOT EXISTS meal_shopping_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    daily_meal_id UUID NOT NULL REFERENCES daily_meals(id) ON DELETE CASCADE,
    item_name VARCHAR(255) NOT NULL,
    quantity VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create shopping_list table
CREATE TABLE IF NOT EXISTS shopping_list (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_name VARCHAR(255) NOT NULL,
    quantity VARCHAR(100) NOT NULL,
    is_checked BOOLEAN DEFAULT false,
    source VARCHAR(50) NOT NULL DEFAULT 'manual',
    source_meal_id UUID REFERENCES daily_meals(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    checked_at TIMESTAMP
);

-- Create meal_generation_history table
CREATE TABLE IF NOT EXISTS meal_generation_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    meal_name VARCHAR(255) NOT NULL,
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_proposals_user_date ON weekly_proposals(user_id, week_start_date DESC);
CREATE INDEX IF NOT EXISTS idx_daily_meals_proposal ON daily_meals(proposal_id);
CREATE INDEX IF NOT EXISTS idx_shopping_list_user ON shopping_list(user_id, is_checked);
CREATE INDEX IF NOT EXISTS idx_history_user ON meal_generation_history(user_id, generated_at DESC);
CREATE INDEX IF NOT EXISTS idx_meal_shopping_items_meal ON meal_shopping_items(daily_meal_id);
