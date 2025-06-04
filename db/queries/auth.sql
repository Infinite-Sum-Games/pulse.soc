-- name: CheckUserExistQuery :one
SELECT EXISTS(
    SELECT 1
    FROM user_account
    WHERE
      ghUsername = $1
  );

-- name: RetriveExistingUserQuery :one
SELECT
  email,
  ghUsername
FROM
  user_account
WHERE
  ghUsername = $1;

-- name: AddRefreshTokenQuery :one
UPDATE user_account
SET
  refresh_token = $1,
  updated_at = NOW()
WHERE
  ghUsername = $2
  AND status = true
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
    first_name,
    middle_name,
    last_name,
    email,
    ghUsername,
    otp,
    expiry_at
  )
VALUES ($1, $2, $3, $4, $5, $6, NOW() + INTERVAL '7 minutes')
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
  first_name, middle_name, last_name, email, ghUsername;

-- name: CreateUserAccountQuery :one
INSERT INTO
  user_account
  (
    first_name,
    middle_name,
    last_name,
    email,
    ghUsername
  )
VALUES ($1, $2, $3, $4, $5)
RETURNING
  ghUsername;
