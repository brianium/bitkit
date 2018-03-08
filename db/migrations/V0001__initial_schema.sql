CREATE SCHEMA bitkit;

-- What precision and scale (if needed) should these numeric types take?
-- What indexes did we want?
CREATE TABLE bitkit.transactions (
    id VARCHAR(64) PRIMARY KEY,
    fee_rate REAL NOT NULL,
    weight INTEGER NOT NULL
);

CREATE INDEX fee_rate_idx ON bitkit.transactions (fee_rate);
