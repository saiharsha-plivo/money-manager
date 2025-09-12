CREATE TABLE IF NOT EXISTS records (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    amount INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    description TEXT,
    type_id INT NOT NULL,
    currency_id INT NOT NULL,
    version INT NOT NULL DEFAULT 1
);

ALTER TABLE records ADD CONSTRAINT fk_records_types FOREIGN KEY (type_id) REFERENCES types (id);
ALTER TABLE records ADD CONSTRAINT fk_records_currencies FOREIGN KEY (currency_id) REFERENCES currencies (id);