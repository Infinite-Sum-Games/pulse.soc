-- name: FetchProfileQuery :one
SELECT
  ua.email,
  ua.ghUsername,
  ua.first_name,
  ua.middle_name,
  ua.last_name,
  ua.bounty,
  COUNT(DISTINCT s.id) as solution_count,
  COUNT(DISTINCT ic.id) as active_claims
FROM 
  user_account ua
LEFT JOIN solutions s ON s.ghUsername = ua.ghUsername
LEFT JOIN issue_claims ic ON ic.ghUsername = ua.ghUsername
LEFT JOIN issues i ON i.id = ic.issue_id AND i.resolved = false AND ic.elapsed_on > NOW()
WHERE
  ua.status = true
  AND ua.ghUsername = $1
GROUP BY
  ua.email,
  ua.ghUsername,
  ua.first_name,
  ua.middle_name,
  ua.last_name,
  ua.bounty;
  
-- name: FetchBadgesQuery :many
SELECT
  badge_name, 
  awarded_at
FROM 
  badge_dispatch
WHERE
  ghUsername = $1;