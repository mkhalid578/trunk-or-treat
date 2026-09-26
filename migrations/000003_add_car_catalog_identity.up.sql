CREATE UNIQUE INDEX cars_catalog_identity_idx
ON cars (LOWER(make), LOWER(model), model_year, LOWER(COALESCE(trim, '')));
