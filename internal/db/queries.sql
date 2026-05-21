-- name: CreateApplication :one
INSERT INTO applications (id, name, description, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING *;

-- name: GetApplication :one
SELECT * FROM applications WHERE id = $1;

-- name: ListApplications :many
SELECT * FROM applications ORDER BY created_at DESC;

-- name: UpdateApplication :one
UPDATE applications
SET name = $2, description = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteApplication :exec
DELETE FROM applications WHERE id = $1;

-- name: CreateFeatureFlag :one
INSERT INTO feature_flags (id, app_id, name, description, enabled, rollout_percentage, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
RETURNING *;

-- name: GetFeatureFlag :one
SELECT * FROM feature_flags WHERE id = $1;

-- name: ListFeatureFlagsByApp :many
SELECT * FROM feature_flags WHERE app_id = $1 ORDER BY name;

-- name: ListAllFeatureFlags :many
SELECT * FROM feature_flags ORDER BY app_id, name;

-- name: UpdateFeatureFlag :one
UPDATE feature_flags
SET name = $2, description = $3, enabled = $4, rollout_percentage = $5, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateFeatureFlagEnabled :one
UPDATE feature_flags
SET enabled = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateFeatureFlagRollout :one
UPDATE feature_flags
SET rollout_percentage = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteFeatureFlag :exec
DELETE FROM feature_flags WHERE id = $1;
