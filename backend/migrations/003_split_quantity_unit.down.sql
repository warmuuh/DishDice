-- Rollback: remove unit columns
DROP INDEX IF EXISTS idx_shopping_list_item_unit;
ALTER TABLE shopping_list DROP COLUMN IF EXISTS unit;
ALTER TABLE meal_shopping_items DROP COLUMN IF EXISTS unit;
