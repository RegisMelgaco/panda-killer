CREATE TABLE IF NOT EXISTS transaction (
   transaction_id SERIAL PRIMARY KEY,
   origin_account SERIAL REFERENCES account (account_id),
   destination_account SERIAL REFERENCES account (account_id),
   amount DOUBLE PRECISION NOT NULL,
   created_at TIMESTAMP WITH TIME ZONE NOT NULL
);
