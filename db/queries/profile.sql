-- name: FetchProfileQuery :one
SELECT
  email,
  ghUsername,
  first_name,
  middle_name,
  last_name,
  bounty
FROM 
  user_account
WHERE
  status = true
  AND ghUsername = $1;
  
-- name: FetchBadgesQuery :many
SELECT
  badge_id, 
  awarded_at
FROM 
  badge_dispatch
WHERE
  ghUsername = $1;

-- name: FetchStatsQuery :one

