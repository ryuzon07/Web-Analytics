-- name: CreateEvent :one
INSERT INTO events (
  site_id, event_type, path, user_id, "timestamp"
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING id;

-- name: GetSiteStats :one
SELECT
    COALESCE(COUNT(id), 0)::bigint AS total_views,
    COALESCE(COUNT(DISTINCT user_id), 0)::bigint AS unique_users
FROM events
WHERE
    site_id = $1
    AND event_type = 'page_view'
    AND "timestamp" >= $2
    AND "timestamp" < $3;

-- name: GetTopPaths :many
SELECT
    path,
    COUNT(id) AS views
FROM events
WHERE
    site_id = $1
    AND event_type = 'page_view'
    AND "timestamp" >= $2
    AND "timestamp" < $3
GROUP BY path
ORDER BY views DESC
LIMIT 10;