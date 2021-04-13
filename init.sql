CREATE TABLE IF NOT EXISTS locations (
    location_id SERIAL NOT NULL,
    truck_id INT NOT NULL,
    latitude FLOAT(32) NOT NULL,
    longitude FLOAT(32) NOT NULL,
    location_time TIMESTAMP NOT NULL,
    PRIMARY KEY(location_id)
);