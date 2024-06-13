-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    user_id INT NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    username VARCHAR(255) NOT NULL UNIQUE,
    hash VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS creds (
    cred_id INT NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id INT NOT NULL REFERENCES users(user_id),
    user_site VARCHAR(255) NOT NULL UNIQUE,
    user_login VARCHAR(255) NOT NULL,
    user_password VARCHAR(255) NOT NULL,
    metadata  VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS creds;
-- +goose StatementEnd
