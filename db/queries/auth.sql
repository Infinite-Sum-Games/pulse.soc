-- name: CheckUserExistQuery :one
SELECT
  ghUsername,
  email
FROM 
  user_account
WHERE
  status = true
  AND ghUsername = $1;

-- name: AddRefreshTokenQuery :one
UPDATE user_account
SET
  refresh_token = $1
WHERE
  ghUsername = $2
  AND status = true
  AND updated_at = NOW()
RETURNING
  email,
  ghUsername,
  refresh_token,
  bounty;

-- name: CheckRefreshTokenQuery :one
SELECT
  ghUsername, email
FROM
  user_account
WHERE
  refresh_token = $1
  AND email = $2
  AND status = true;

-- name: CheckForExistingOtpQuery :one
SELECT
  email,
  otp
FROM
  user_onboarding 
WHERE
  ghUsername = $1
  AND expiry_at >= NOW() + INTERVAL '1 minute';

-- name: BeginUserRegistrationQuery :one
INSERT INTO 
  user_onboarding
  (
    email,
    ghUsername,
    otp,
    expiry_at
  )
VALUES ($1, $2, $3, NOW() + INTERVAL '5 minutes')
RETURNING
  email, otp;

-- name: VerifyOtpQuery :one
DELETE FROM 
  user_onboarding
WHERE
  ghUsername = $1
  AND otp = $2
  AND expiry_at > NOW()
RETURNING
  email, ghUsername;

-- name: CreateUserAccountQuery :one
INSERT INTO
  user_account
  (
    email,
    ghUsername
  )
VALUES ($1, $2)
RETURNING
  ghUsername;
