-- name: InsertAccount :one
INSERT INTO account(name, cpf, secret, balance, created_at) values($1, $2, $3, $4, $5) RETURNING account_id;