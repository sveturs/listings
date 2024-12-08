CREATE TABLE cars (
    id SERIAL PRIMARY KEY,
    make VARCHAR(50),
    model VARCHAR(50),
    year INT,
    price_per_day NUMERIC(10, 2),
    availability BOOLEAN DEFAULT TRUE,
    location VARCHAR(100)
);