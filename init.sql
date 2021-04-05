CREATE TABLE IF NOT EXISTS drivers (
    driver_id SERIAL PRIMARY KEY,
    driver_name TEXT NOT NULL,
    driver_surname TEXT NOT NULL,
    driver_age INT,
    driver_email TEXT,
    driver_address TEXT,
    driver_phone TEXT
);

CREATE TABLE IF NOT EXISTS locations (
    location_id SERIAL NOT NULL,
    driver_id INT NOT NULL,
    latitude FLOAT(32) NOT NULL,
    longitude FLOAT(32) NOT NULL,
    PRIMARY KEY(location_id),
    FOREIGN KEY (driver_id) REFERENCES drivers(driver_id)
);