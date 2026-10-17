CREATE UNIQUE INDEX cars_make_model_year_unique_idx
ON cars (LOWER(make), LOWER(model), model_year);
