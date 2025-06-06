-- name: FetchParticipantListQuery :many
SELECT
  CASE
    WHEN middle_name IS NULL THEN first_name || ' ' || last_name
    ELSE first_name || ' ' || middle_name || ' ' || last_name
  END as full_name,
  ghUsername as github_username,
  0 as bounty,
  0 as solutions
FROM
  user_account
WHERE
  status = true;