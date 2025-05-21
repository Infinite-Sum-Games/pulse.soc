-- name: FetchProfileQuery :one
SELECT
  email,
  ghUsername,
  bounty
FROM 
  user_account
WHERE
  status = true
  AND ghUsername = $1
  AND refresh_token IS NOT NULL;
  
-- name: FetchBadgesQuery :many
SELECT
  badge_id, 
  awarded_at
FROM 
  badge_dispatch
WHERE
  ghUsername = $1;

