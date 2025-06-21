CREATE TABLE IF NOT EXISTS goauth.users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			email_verified BOOLEAN DEFAULT false NOT NULL,
			verification_token VARCHAR(300),
			verification_expires_at TIMESTAMP
		);
