-- name: CheckRefreshTokenQuery :one
SELECT
  ghUsername, email
FROM
  user_account
WHERE
  refresh_token = $1
  AND email = $2
  AND status = true;

-- name: OnboardUserQuery :one

-- name: AddRefreshTokenQuery :one
