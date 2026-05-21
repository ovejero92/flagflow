package db

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FeatureFlag struct {
	ID                uuid.UUID `json:"id"`
	AppID             uuid.UUID `json:"app_id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	Enabled           bool      `json:"enabled"`
	RolloutPercentage int32     `json:"rollout_percentage"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
