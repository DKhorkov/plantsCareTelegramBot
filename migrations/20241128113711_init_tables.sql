-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id          SERIAL PRIMARY KEY,
    telegram_id BIGINT       NOT NULL UNIQUE,
    username    VARCHAR(100) NOT NULL UNIQUE,
    firstname   VARCHAR(100) NOT NULL,
    lastname    VARCHAR(100) NOT NULL,
    is_bot      BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS groups
(
    id                      SERIAL PRIMARY KEY,
    user_id                 INTEGER      NOT NULL,
    title                   VARCHAR(100) NOT NULL,
    description             TEXT,
    last_watering_date      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    next_watering_date      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    watering_interval_hours INTEGER      NOT NULL,
    created_at              TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS plants
(
    id          SERIAL PRIMARY KEY,
    group_id    INTEGER      NOT NULL,
    user_id     INTEGER      NOT NULL,
    title       VARCHAR(100) NOT NULL,
    description TEXT,
    photo       bytea,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS notifications
(
    id         SERIAL PRIMARY KEY,
    group_id   INTEGER   NOT NULL,
    message_id INTEGER   NOT NULL,
    sent_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS temporary
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL,
    step       INTEGER NOT NULL,
    message_id BIGINT,
    data       bytea,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS temporary;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS plants;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
