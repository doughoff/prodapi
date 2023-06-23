BEGIN;

-- Remove the unique constraint from the 'code' column
ALTER TABLE production_orders DROP CONSTRAINT code_unique;

-- Alter the 'code' column to be NOT NULL
ALTER TABLE production_orders ALTER COLUMN code SET NOT NULL;

COMMIT;
