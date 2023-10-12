-- name: GetUsers :many
SELECT id, username, discord, created_at, active_at, permission FROM users;

-- name: SetUserPass :exec
UPDATE users SET password = $1 WHERE id = $2;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByName :one
SELECT * FROM users WHERE username = $1;

-- name: GetBingos :many
SELECT * from bingos WHERE visible = true;

-- name: AddBingo :one
INSERT INTO bingos (name, size) VALUES ($1, $2)
RETURNING *;

-- name: AddBingoField :exec
INSERT INTO bingo_fields (text, bingo_id) VALUES ($1, $2);

-- name: GetBingoFields :many
SELECT * from bingo_fields WHERE bingo_id = $1;

-- name: GetBingoUsers :many
SELECT bingo_users.*, users.username from bingo_users
JOIN users ON bingo_users.user_id = users.id
WHERE bingo_users.bingo_id = $1 AND users.permission IN ('Moderator', 'User');

-- name: JoinBingo :exec
INSERT INTO bingo_users (user_id, bingo_id)
VALUES ($1, $2);

-- name: GetBingoUserFields :many
SELECT * FROM bingo_user_fields
WHERE bingo_id = $1 AND user_id = $2;

-- name: GetBingoUserField :one
SELECT * FROM bingo_user_fields WHERE user_id = $1 AND bingo_id = $2 AND bingo_field_id = $3;

-- name: SetBingoUserFieldStatus :exec
INSERT INTO bingo_user_fields(user_id, bingo_id, bingo_field_id, done_at, confirmed_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id, bingo_id, bingo_field_id) DO UPDATE SET done_at = $4, confirmed_at = $5;

-- name: GetUnconfirmedUserFields :many
SELECT bingo_user_fields.*, bingo_fields.text as field_label, users.username FROM bingo_user_fields
JOIN bingo_fields ON bingo_fields.id = bingo_user_fields.bingo_field_id
JOIN users ON bingo_user_fields.user_id = users.id
WHERE confirmed_at IS NULL;

-- name: InsertPOEUser :one
INSERT INTO users (poe_id, username, permission)
VALUES ($1, $2, $3)
ON CONFLICT (poe_id) DO UPDATE set username = $2
RETURNING *;

-- name: SetUserToken :exec
INSERT INTO user_tokens (user_id, refresh_token, access_token, expires_at, scope)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id) DO UPDATE SET refresh_token = $2, access_token = $3, expires_at = $4, scope = $5;

-- name: SetUserPermission :exec
UPDATE users
SET permission = $2
WHERE id = $1;
