CREATE TABLE IF NOT EXISTS account(
   account_id serial PRIMARY KEY,
   name VARCHAR (50) NOT NULL,
   cpf VARCHAR (11) NOT NULL,
   secret VARCHAR (255) NOT NULL
      CHECK(secret <> ''),
   balance DOUBLE PRECISION NOT NULL,
   created_at TIMESTAMP WITH TIME ZONE NOT NULL
);