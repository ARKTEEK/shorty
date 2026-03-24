-- name: CreateLink :execresult
INSERT INTO links (original_url, short_code, user_id)
VALUES (?, ?, ?);

-- name: GetLinkByCode :one
SELECT * FROM links
WHERE short_code = ? LIMIT 1;

-- name: GetLinksByUser :many
SELECT * FROM links
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: CreateUser :execresult
INSERT INTO users (email, password)
VALUES (?, ?);

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ? LIMIT 1;
