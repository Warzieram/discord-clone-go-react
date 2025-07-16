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

/*
CREATE TABLE IF NOT EXISTS chatroom.rooms (
  id SERIAL PRIMARY KEY,
);
*/

CREATE TABLE IF NOT EXISTS chatroom.messages (
  id SERIAL PRIMARY KEY,
  content VARCHAR(500),
  sender_id int,
  room_id int,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  --CONSTRAINT fk_room_id FOREIGN KEY (room_id)
  --REFERENCES chatroom.rooms(id),
  CONSTRAINT fk_sender_id FOREIGN KEY (sender_id)
  REFERENCES chatroom.users(id)
);
/*
CREATE TABLE IF NOT EXISTS chatroom.has_user (
  message_id int,
  room_id int,
  CONSTRAINT fk_message_id FOREIGN KEY (sender_id)
  REFERENCES chatroom.users(id),
  CONSTRAINT fk_room_id FOREIGN KEY (room_id)
  REFERENCES chatroom.rooms(id),
)
  */
