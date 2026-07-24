-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (
    username, 
    email, 
    secret_code
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: UpdateVerifyEmail :one
UPDATE verify_emails
SET is_used = TRUE
WHERE id = sqlc.arg(id)
  AND secret_code = sqlc.arg(secret_code)
  AND is_used = FALSE
  AND expires_at > now()
RETURNING *;
