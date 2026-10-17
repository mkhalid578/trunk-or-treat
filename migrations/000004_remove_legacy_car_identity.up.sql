DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'cars_make_model_year_unique_idx'
          AND conrelid = 'cars'::regclass
    ) THEN
        ALTER TABLE cars
        DROP CONSTRAINT cars_make_model_year_unique_idx;
    END IF;
END
$$;

DROP INDEX IF EXISTS cars_make_model_year_unique_idx;
