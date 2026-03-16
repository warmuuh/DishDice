-- Add unit column to meal_shopping_items
ALTER TABLE meal_shopping_items ADD COLUMN unit VARCHAR(50);

-- Add unit column to shopping_list
ALTER TABLE shopping_list ADD COLUMN unit VARCHAR(50);

-- Try to parse existing quantities and split them
-- This is best-effort for existing data
UPDATE meal_shopping_items
SET unit = CASE
    WHEN quantity LIKE '%g' THEN 'g'
    WHEN quantity LIKE '%kg' THEN 'kg'
    WHEN quantity LIKE '%ml' THEN 'ml'
    WHEN quantity LIKE '%l' THEN 'l'
    WHEN quantity LIKE '%TL%' OR quantity LIKE '%Teelöffel%' THEN 'TL'
    WHEN quantity LIKE '%EL%' OR quantity LIKE '%Esslöffel%' THEN 'EL'
    WHEN quantity LIKE '%Prise%' THEN 'Prise'
    WHEN quantity LIKE '%Stück%' OR quantity LIKE '%Stk%' THEN 'Stück'
    ELSE ''
END;

UPDATE shopping_list
SET unit = CASE
    WHEN quantity LIKE '%g' THEN 'g'
    WHEN quantity LIKE '%kg' THEN 'kg'
    WHEN quantity LIKE '%ml' THEN 'ml'
    WHEN quantity LIKE '%l' THEN 'l'
    WHEN quantity LIKE '%TL%' OR quantity LIKE '%Teelöffel%' THEN 'TL'
    WHEN quantity LIKE '%EL%' OR quantity LIKE '%Esslöffel%' THEN 'EL'
    WHEN quantity LIKE '%Prise%' THEN 'Prise'
    WHEN quantity LIKE '%Stück%' OR quantity LIKE '%Stk%' THEN 'Stück'
    ELSE ''
END;

-- Update quantity column to contain only numbers
-- This is approximate and will need manual review for existing data
UPDATE meal_shopping_items
SET quantity = TRIM(REGEXP_REPLACE(quantity, '[^0-9.,]', '', 'g'))
WHERE quantity IS NOT NULL;

UPDATE shopping_list
SET quantity = TRIM(REGEXP_REPLACE(quantity, '[^0-9.,]', '', 'g'))
WHERE quantity IS NOT NULL;

-- Create index for faster shopping list merging queries
CREATE INDEX IF NOT EXISTS idx_shopping_list_item_unit ON shopping_list(user_id, item_name, unit);
