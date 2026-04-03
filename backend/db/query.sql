-- name: GetApplication :one
SELECT * FROM applications
WHERE id = ? LIMIT 1;

-- name: GetApplicationByEmail :one
SELECT * FROM applications
WHERE email = ? LIMIT 1;

-- name: CreateApplication :exec
INSERT INTO applications (id, email, code, expires_at)
VALUES (?, ?, ?, ?);

-- name: UpdateApplication :exec
UPDATE applications
SET email = ?, code = ?, expires_at = ?
WHERE id = ?;

-- name: DeleteApplication :exec
DELETE FROM applications
WHERE id = ?;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ? LIMIT 1;
