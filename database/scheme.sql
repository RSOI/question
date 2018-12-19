SET SYNCHRONOUS_COMMIT = 'off';
CREATE EXTENSION IF NOT EXISTS CITEXT;
CREATE SCHEMA IF NOT EXISTS question;

DROP TABLE IF EXISTS question.question;
DROP TABLE IF EXISTS question.services;

CREATE TABLE question.question (
	id SERIAL PRIMARY KEY,
	title CITEXT NOT NULL,
	content CITEXT NULL,
	author_id INTEGER NOT NULL,
	author_nickname CITEXT NOT NULL,
	has_best BOOLEAN DEFAULT FALSE,
	created TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS author_id_index ON question.question (author_id);
CREATE INDEX IF NOT EXISTS has_best_index ON question.question (has_best);
CREATE INDEX IF NOT EXISTS author_id__has_best_index ON question.question (author_id, has_best);

CREATE TABLE question.services (
	id SERIAL PRIMARY KEY,
	request CITEXT NOT NULL,
	request_time TIMESTAMPTZ DEFAULT NOW(),
	response_status INTEGER NOT NULL,
	response_error_text CITEXT NULL
);
