WITH grouped_cars AS (
    SELECT
        LOWER(make) AS make_key,
        LOWER(model) AS model_key,
        model_year,
        MIN(id) AS keep_id,
        STRING_AGG(
            DISTINCT NULLIF(BTRIM(trim), ''),
            ', ' ORDER BY NULLIF(BTRIM(trim), '')
        ) AS combined_trims,
        MAX(cargo_vol_seats_folded_cu_ft) AS max_cargo_seats_folded,
        MAX(cargo_vol_behind_2nd_row_cu_ft) AS max_cargo_behind_2nd_row,
        CASE
            WHEN COUNT(DISTINCT body_style) > 1 THEN 'other'
            ELSE MIN(body_style)
        END AS combined_body_style,
        CASE
            WHEN COUNT(DISTINCT powertrain) > 1 THEN 'other'
            ELSE MIN(powertrain)
        END AS combined_powertrain
    FROM cars
    GROUP BY LOWER(make), LOWER(model), model_year
)
UPDATE cars AS car
SET trim = grouped.combined_trims,
    body_style = grouped.combined_body_style,
    powertrain = grouped.combined_powertrain,
    cargo_vol_seats_folded_cu_ft = COALESCE(car.cargo_vol_seats_folded_cu_ft, grouped.max_cargo_seats_folded),
    cargo_vol_behind_2nd_row_cu_ft = COALESCE(car.cargo_vol_behind_2nd_row_cu_ft, grouped.max_cargo_behind_2nd_row)
FROM grouped_cars AS grouped
WHERE car.id = grouped.keep_id;

WITH ranked_cars AS (
    SELECT
        id,
        ROW_NUMBER() OVER (
            PARTITION BY LOWER(make), LOWER(model), model_year
            ORDER BY id
        ) AS row_number
    FROM cars
)
DELETE FROM cars
WHERE id IN (
    SELECT id
    FROM ranked_cars
    WHERE row_number > 1
);

CREATE UNIQUE INDEX cars_make_model_year_unique_idx
ON cars (LOWER(make), LOWER(model), model_year);
