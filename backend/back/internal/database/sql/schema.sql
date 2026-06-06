CREATE SCHEMA IF NOT EXISTS ${BLUEPRINT_DB_SCHEMA};

CREATE TABLE IF NOT EXISTS ${BLUEPRINT_DB_SCHEMA}.users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
      username VARCHAR(50) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			email_verified BOOLEAN DEFAULT false NOT NULL,
			verification_token VARCHAR(300),
			verification_expires_at TIMESTAMP
		);

/*
CREATE TABLE IF NOT EXISTS ${BLUEPRINT_DB_SCHEMA}.rooms (
  id SERIAL PRIMARY KEY,
);
*/

CREATE TABLE IF NOT EXISTS ${BLUEPRINT_DB_SCHEMA}.messages (
  id SERIAL PRIMARY KEY,
  content VARCHAR(500),
  sender_id INT,
  room_id INT,
  deleted BOOLEAN DEFAULT false,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  --CONSTRAINT fk_room_id FOREIGN KEY (room_id)
  --REFERENCES ${BLUEPRINT_DB_SCHEMA}.rooms(id),
  CONSTRAINT fk_sender_id FOREIGN KEY (sender_id)
  REFERENCES ${BLUEPRINT_DB_SCHEMA}.users(id)
);
/*
CREATE TABLE IF NOT EXISTS ${BLUEPRINT_DB_SCHEMA}.has_user (
  message_id int,
  room_id int,
  CONSTRAINT fk_message_id FOREIGN KEY (sender_id)
  REFERENCES chatroom.users(id),
  CONSTRAINT fk_room_id FOREIGN KEY (room_id)
  REFERENCES chatroom.rooms(id),
)
  */
