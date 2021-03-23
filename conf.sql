CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
	id			uuid		DEFAULT uuid_generate_v4 (),
	username	text		NOT NULL,
	password	text		NOT NULL,
	UNIQUE(id),
	UNIQUE(username)
);

CREATE TABLE posts (
	id 			uuid		DEFAULT uuid_generate_v4 (),
	author		text		NOT NULL references users(username),
	time_stamp	timestamp	DEFAULT CURRENT_TIMESTAMP,
	title		text		NOT NULL,
	context 	text		NOT NULL
);