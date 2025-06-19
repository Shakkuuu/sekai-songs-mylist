-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY id;

-- name: InsertUser :one
INSERT INTO users (id, email, password, is_verified, verify_token, token_expires_at, is_admin, created_at, updated_at, deleted_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: ExistsUserByEmail :one
SELECT EXISTS (
  SELECT 1 FROM users WHERE email = $1
) AS exists;

-- name: ExistsUserByID :one
SELECT EXISTS (
  SELECT 1 FROM users WHERE id = $1
) AS exists;

-- name: IsVerifiedByID :one
SELECT is_verified FROM users WHERE id = $1;

-- name: IsAdminByID :one
SELECT is_admin FROM users WHERE id = $1;

-- name: UpdateUserEmail :exec
UPDATE users
SET email = $1,
    updated_at = $2
WHERE id = $3;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = $1,
    updated_at = $2
WHERE id = $3;

-- name: UpdateUserIsVerified :exec
UPDATE users
SET is_verified = $1,
    updated_at = $2
WHERE id = $3;

-- name: UpdateUserVerifyToken :exec
UPDATE users
SET verify_token = $1,
    updated_at = $2
WHERE id = $3;

-- name: UpdateUserTokenExpiresAt :exec
UPDATE users
SET token_expires_at = $1,
    updated_at = $2
WHERE id = $3;

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = $1,
    updated_at = $2
WHERE id = $3;
