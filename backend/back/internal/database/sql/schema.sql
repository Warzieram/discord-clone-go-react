CREATE SCHEMA IF NOT EXISTS chatroom;

CREATE TABLE IF NOT EXISTS chatroom.users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
      username VARCHAR(50) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			email_verified BOOLEAN DEFAULT false NOT NULL,
			verification_token VARCHAR(300),
			verification_expires_at TIMESTAMP
		);


CREATE TABLE IF NOT EXISTS chatroom.messages (
  id SERIAL PRIMARY KEY,
  content VARCHAR(500),
  sender_id int,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  constraint fk_sender_id foreign key (sender_id)
  REFERENCES chatroom.users(id)
);

