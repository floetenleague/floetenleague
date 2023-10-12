CREATE TYPE user_permission AS ENUM ('Banned', 'Unverified', 'User', 'Moderator');

CREATE TABLE users (
    id         BIGSERIAL PRIMARY KEY,
    poe_id     TEXT NOT NULL,
    username   VARCHAR(100) NOT NULL,
    password   VARCHAR(100),
    discord    VARCHAR(100),
    permission user_permission NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    active_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(username),
    UNIQUE(poe_id)
);

CREATE TABLE user_tokens (
    user_id       BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token  TEXT NOT NULL,
    access_token  TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    scope TEXT NOT NULL,

    UNIQUE(user_id)
);

CREATE TABLE bingos (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    size       INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bingo_users (
    user_id       BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bingo_id      BIGINT NOT NULL REFERENCES bingos(id) ON DELETE CASCADE,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, bingo_id)
);

CREATE TABLE bingo_fields (
    id         BIGSERIAL PRIMARY KEY,
    bingo_id   BIGINT NOT NULL REFERENCES bingos(id) ON DELETE CASCADE,
    text       TEXT NOT NULL
);


CREATE TABLE bingo_user_fields (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bingo_id BIGINT NOT NULL REFERENCES bingos(id) ON DELETE CASCADE,
    bingo_field_id BIGINT NOT NULL REFERENCES bingo_fields(id) ON DELETE CASCADE,
    done_at TIMESTAMP,
    confirmed_at TIMESTAMP,

    UNIQUE(user_id, bingo_id, bingo_field_id)
);
