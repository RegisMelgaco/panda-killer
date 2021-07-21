-- name: InsertAccount :one
INSERT INTO account(name, cpf, secret, balance, created_at) values($1, $2, $3, $4, $5) RETURNING account_id;

-- name: SelectAccountByCPF :one
SELECT account_id, name, cpf, secret, balance, created_at
			FROM account
			WHERE cpf = $1
			FETCH FIRST ROW ONLY;

-- name: ListAccounts :many
SELECT account_id, name, cpf, secret, balance, created_at FROM account;

-- name: GetAccount :one
SELECT account_id, name, cpf, secret, balance, created_at FROM account WHERE account_id = $1 FETCH FIRST ROW ONLY;
