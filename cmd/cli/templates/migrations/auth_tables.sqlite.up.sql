PRAGMA recursive_triggers = ON;
PRAGMA foreign_keys = ON;

drop table if exists users;

CREATE TABLE users (
                       id INTEGER PRIMARY KEY AUTOINCREMENT,
                       first_name character varying(255) NOT NULL,
                       last_name character varying(255) NOT NULL,
                       user_active integer NOT NULL DEFAULT 0,
                       email character varying(255) NOT NULL UNIQUE,
                       password character varying(60) NOT NULL,
                       created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP

);


CREATE TRIGGER set_timestamp_on_users
    BEFORE UPDATE ON users
    FOR EACH ROW
BEGIN
UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE created_at < CURRENT_TIMESTAMP;
END;

drop table if exists remember_tokens;

CREATE TABLE `remember_tokens` (
                                   id INTEGER PRIMARY KEY AUTOINCREMENT,
                                   user_id integer NOT NULL REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,
    remember_token character varying(100) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE TRIGGER set_timestamp_on_remember_tokens
    BEFORE UPDATE ON remember_tokens
    FOR EACH ROW
BEGIN
UPDATE remember_tokens SET updated_at = CURRENT_TIMESTAMP WHERE created_at < CURRENT_TIMESTAMP;
END;

drop table if exists tokens;

CREATE TABLE `tokens` (
                          id INTEGER PRIMARY KEY AUTOINCREMENT,
                          user_id integer NOT NULL REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,
    first_name character varying(255) NOT NULL,
    email character varying(255) NOT NULL UNIQUE,
    tokens character varying(255) NOT NULL,
    token character varying(255) NOT NULL,
    token_hash bytea NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL
    );

CREATE TRIGGER set_timestamp_on_tokens
    BEFORE UPDATE ON tokens
    FOR EACH ROW
BEGIN
UPDATE tokens SET updated_at = CURRENT_TIMESTAMP WHERE created_at < CURRENT_TIMESTAMP;
END;