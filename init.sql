CREATE TABLE IF NOT EXISTS drivers (
    driver_id SERIAL PRIMARY KEY,
    driver_name TEXT,
    driver_surname TEXT,
    driver_age INT,
    driver_email TEXT UNIQUE NOT NULL,
    driver_address TEXT,
    driver_phone TEXT
);