package db

import (
	"context"

	"github.com/google/uuid"
)

const createApplication = `-- name: CreateApplication :one
INSERT INTO applications (id, name, description, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING id, name, description, created_at, updated_at`

type CreateApplicationParams struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func (q *Queries) CreateApplication(ctx context.Context, arg CreateApplicationParams) (Application, error) {
	row := q.db.QueryRow(ctx, createApplication, arg.ID, arg.Name, arg.Description)
	var i Application
	err := row.Scan(&i.ID, &i.Name, &i.Description, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const getApplication = `-- name: GetApplication :one
SELECT id, name, description, created_at, updated_at FROM applications WHERE id = $1`

func (q *Queries) GetApplication(ctx context.Context, id uuid.UUID) (Application, error) {
	row := q.db.QueryRow(ctx, getApplication, id)
	var i Application
	err := row.Scan(&i.ID, &i.Name, &i.Description, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const listApplications = `-- name: ListApplications :many
SELECT id, name, description, created_at, updated_at FROM applications ORDER BY created_at DESC`

func (q *Queries) ListApplications(ctx context.Context) ([]Application, error) {
	rows, err := q.db.Query(ctx, listApplications)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Application
	for rows.Next() {
		var i Application
		if err := rows.Scan(&i.ID, &i.Name, &i.Description, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const updateApplication = `-- name: UpdateApplication :one
UPDATE applications
SET name = $2, description = $3, updated_at = NOW()
WHERE id = $1
RETURNING id, name, description, created_at, updated_at`

type UpdateApplicationParams struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func (q *Queries) UpdateApplication(ctx context.Context, arg UpdateApplicationParams) (Application, error) {
	row := q.db.QueryRow(ctx, updateApplication, arg.ID, arg.Name, arg.Description)
	var i Application
	err := row.Scan(&i.ID, &i.Name, &i.Description, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const deleteApplication = `-- name: DeleteApplication :exec
DELETE FROM applications WHERE id = $1`

func (q *Queries) DeleteApplication(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteApplication, id)
	return err
}

const createFeatureFlag = `-- name: CreateFeatureFlag :one
INSERT INTO feature_flags (id, app_id, name, description, enabled, rollout_percentage, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
RETURNING id, app_id, name, description, enabled, rollout_percentage, created_at, updated_at`

type CreateFeatureFlagParams struct {
	ID                uuid.UUID `json:"id"`
	AppID             uuid.UUID `json:"app_id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	Enabled           bool      `json:"enabled"`
	RolloutPercentage int32     `json:"rollout_percentage"`
}

func (q *Queries) CreateFeatureFlag(ctx context.Context, arg CreateFeatureFlagParams) (FeatureFlag, error) {
	row := q.db.QueryRow(ctx, createFeatureFlag,
		arg.ID, arg.AppID, arg.Name, arg.Description, arg.Enabled, arg.RolloutPercentage)
	var i FeatureFlag
	err := row.Scan(&i.ID, &i.AppID, &i.Name, &i.Description, &i.Enabled, &i.RolloutPercentage, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const getFeatureFlag = `-- name: GetFeatureFlag :one
SELECT id, app_id, name, description, enabled, rollout_percentage, created_at, updated_at FROM feature_flags WHERE id = $1`

func (q *Queries) GetFeatureFlag(ctx context.Context, id uuid.UUID) (FeatureFlag, error) {
	row := q.db.QueryRow(ctx, getFeatureFlag, id)
	var i FeatureFlag
	err := row.Scan(&i.ID, &i.AppID, &i.Name, &i.Description, &i.Enabled, &i.RolloutPercentage, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const listFeatureFlagsByApp = `-- name: ListFeatureFlagsByApp :many
SELECT id, app_id, name, description, enabled, rollout_percentage, created_at, updated_at FROM feature_flags WHERE app_id = $1 ORDER BY name`

func (q *Queries) ListFeatureFlagsByApp(ctx context.Context, appID uuid.UUID) ([]FeatureFlag, error) {
	rows, err := q.db.Query(ctx, listFeatureFlagsByApp, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeatureFlag
	for rows.Next() {
		var i FeatureFlag
		if err := rows.Scan(&i.ID, &i.AppID, &i.Name, &i.Description, &i.Enabled, &i.RolloutPercentage, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const listAllFeatureFlags = `-- name: ListAllFeatureFlags :many
SELECT id, app_id, name, description, enabled, rollout_percentage, created_at, updated_at FROM feature_flags ORDER BY app_id, name`

func (q *Queries) ListAllFeatureFlags(ctx context.Context) ([]FeatureFlag, error) {
	rows, err := q.db.Query(ctx, listAllFeatureFlags)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeatureFlag
	for rows.Next() {
		var i FeatureFlag
		if err := rows.Scan(&i.ID, &i.AppID, &i.Name, &i.Description, &i.Enabled, &i.RolloutPercentage, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const updateFeatureFlag = `-- name: UpdateFeatureFlag :one
UPDATE feature_flags
SET name = $2, description = $3, enabled = $4, rollout_percentage = $5, updated_at = NOW()
WHERE id = $1
RETURNING id, app_id, name, description, enabled, rollout_percentage, created_at, updated_at`

type UpdateFeatureFlagParams struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	Enabled           bool      `json:"enabled"`
	RolloutPercentage int32     `json:"rollout_percentage"`
}

func (q *Queries) UpdateFeatureFlag(ctx context.Context, arg UpdateFeatureFlagParams) (FeatureFlag, error) {
	row := q.db.QueryRow(ctx, updateFeatureFlag,
		arg.ID, arg.Name, arg.Description, arg.Enabled, arg.RolloutPercentage)
	var i FeatureFlag
	err := row.Scan(&i.ID, &i.AppID, &i.Name, &i.Description, &i.Enabled, &i.RolloutPercentage, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const updateFeatureFlagEnabled = `-- name: UpdateFeatureFlagEnabled :one
UPDATE feature_flags SET enabled = $2, updated_at = NOW() WHERE id = $1
RETURNING id, app_id, name, description, enabled, rollout_percentage, created_at, updated_at`

func (q *Queries) UpdateFeatureFlagEnabled(ctx context.Context, id uuid.UUID, enabled bool) (FeatureFlag, error) {
	row := q.db.QueryRow(ctx, updateFeatureFlagEnabled, id, enabled)
	var i FeatureFlag
	err := row.Scan(&i.ID, &i.AppID, &i.Name, &i.Description, &i.Enabled, &i.RolloutPercentage, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const updateFeatureFlagRollout = `-- name: UpdateFeatureFlagRollout :one
UPDATE feature_flags SET rollout_percentage = $2, updated_at = NOW() WHERE id = $1
RETURNING id, app_id, name, description, enabled, rollout_percentage, created_at, updated_at`

func (q *Queries) UpdateFeatureFlagRollout(ctx context.Context, id uuid.UUID, rollout int32) (FeatureFlag, error) {
	row := q.db.QueryRow(ctx, updateFeatureFlagRollout, id, rollout)
	var i FeatureFlag
	err := row.Scan(&i.ID, &i.AppID, &i.Name, &i.Description, &i.Enabled, &i.RolloutPercentage, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const deleteFeatureFlag = `-- name: DeleteFeatureFlag :exec
DELETE FROM feature_flags WHERE id = $1`

func (q *Queries) DeleteFeatureFlag(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteFeatureFlag, id)
	return err
}
