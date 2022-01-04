package postgres

const (
	schema = `
		DROP TABLE IF EXISTS transactions;
		DROP TABLE IF EXISTS users;
		
		CREATE TABLE IF NOT EXISTS users 
		(
			id  	UUID PRIMARY KEY,
			balance NUMERIC(10, 2) DEFAULT 0 CHECK (balance >= 0)
		);

		CREATE TABLE IF NOT EXISTS transactions
		(
			id 		   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			to_id      UUID REFERENCES users(id),
			from_id    UUID REFERENCES users(id),
			money      NUMERIC(10, 2) NOT NULL,
			method 	   TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT now()
		);
	`

	sqlCreateUser = `
		INSERT INTO users (id) VALUES ($1);
	`

	sqlUpdateBalance = `
		UPDATE users SET balance = balance + $2
		WHERE id = $1
		RETURNING *;
	`

	sqlGetBalance = `
		SELECT id, balance FROM users 
		WHERE id = $1;
	`

	sqlWriteTransaction = `
		INSERT INTO transactions (to_id, from_id, money, method)
		VALUES ($1, $2, $3, $4);
	`

	sqlGetTransaction = `
		SELECT id, to_id, from_id, money, method, created_at FROM transactions
		WHERE id = $1;
	`

	sqlListOfTransactions = `
		SELECT id, to_id, from_id, money, method, created_at FROM transactions
		WHERE to_id = $1 OR from_id = $1
		%s
		LIMIT $2
		OFFSET ($3 - 1) * $2;
	`
)
