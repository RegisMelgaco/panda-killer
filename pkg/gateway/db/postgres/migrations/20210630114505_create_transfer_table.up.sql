CREATE TABLE IF NOT EXISTS transfer (
   transfer_id SERIAL PRIMARY KEY,
   origin_account SERIAL REFERENCES account (account_id),
   destination_account SERIAL REFERENCES account (account_id),
   amount INTEGER NOT NULL,
   created_at TIMESTAMP WITH TIME ZONE NOT NULL
);
