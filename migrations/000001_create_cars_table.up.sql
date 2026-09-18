CREATE TABLE IF NOT EXISTS cars (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    -- Identity 
    make VARCHAR(50) NOT NULL,
    model VARCHAR(50) NOT NULL,
    model_year SMALLINT NOT NULL CHECK (model_year BETWEEN 1980 and 2100),
    trim TEXT,
    body_style TEXT NOT NULL CHECK (body_style IN ('sedan', 
    'hatchback', 'suv', 'truck', 'coupe', 'convertible', 'wagon', 'van', 'other')),
    powertrain TEXT NOT NULL CHECK (powertrain IN ('gas', 'diesel', 'hybrid', 'electric', 'other')),
    
    -- Cargo volume (cubic feet)
    cargo_vol_cu_ft NUMERIC (5,1) CHECK (cargo_vol_cu_ft >= 0),
    cargo_vol_seats_folded_cu_ft NUMERIC(5,1) CHECK (cargo_vol_seats_folded_cu_ft >= 0),
    cargo_vol_behind_2nd_row_cu_ft NUMERIC(5,1) CHECK (cargo_vol_behind_2nd_row_cu_ft >= 0),
    frunk_vol_cu_ft NUMERIC(5,1) CHECK (frunk_vol_cu_ft >= 0),
    fuel_capacity_gal NUMERIC(5,1) CHECK (fuel_capacity_gal >= 0),
    fuel_economy_mpg NUMERIC(5,1) CHECK (fuel_economy_mpg >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);