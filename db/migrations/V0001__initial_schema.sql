CREATE SCHEMA memcool;

-- What precision and scale (if needed) should these numeric types take?
-- What indexes did we want?
CREATE TABLE memcool.transactions (
    id VARCHAR(64) PRIMARY KEY,
    fee_rate NUMERIC NOT NULL,
    weight NUMERIC NOT NULL
);
