BEGIN;

-- Alter the 'code' column to allow NULL values
ALTER TABLE production_orders ALTER COLUMN code DROP NOT NULL;

-- Add a unique constraint to the 'code' column
ALTER TABLE production_orders ADD CONSTRAINT code_unique UNIQUE (code);

COMMIT;
