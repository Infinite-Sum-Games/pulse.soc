-- name: CheckIfProjectExistsQuery :one
SELECT EXISTS
  (
    SELECT 1 FROM repository
    WHERE id = $1
    LIMIT 1
);

-- name: FetchAllProjectsQuery :many
SELECT 
  id, name, description, url, maintainers, tags, is_internal
FROM 
  repository;

-- name: FetchAllIssuesByProjectIdQuery :many
SELECT
  i.id AS issue_id,
  i.title AS title,
  i.url AS issue_url,
  i.updated_at AS last_update,
  COALESCE(
    JSON_AGG(
      JSON_BUILD_OBJECT(
        'username', c.ghUsername,
        'claimed_on', c.claimed_on,
        'elapsing_on', c.elapsed_on
      ) ORDER BY c.claimed_on
    ) FILTER (WHERE c.id IS NOT NULL),
  '[]'::JSON
  ) AS claimants
FROM issues i
LEFT JOIN
  issue_claims AS c 
  ON c.issue_url = i.url
WHERE 
  i.resolved = false
  AND i.id = $1
GROUP BY
  i.id, 
  i.title, 
  i.url, 
  i.updated_at
ORDER BY
  i.id;
