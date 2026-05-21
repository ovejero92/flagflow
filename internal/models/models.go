package models

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
	ID                 uuid.UUID `json:"id"`
	AppID              uuid.UUID `json:"app_id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	Enabled            bool      `json:"enabled"`
	RolloutPercentage  int       `json:"rollout_percentage"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type CreateApplicationRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateApplicationRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreateFeatureFlagRequest struct {
	Name              string `json:"name" binding:"required"`
	Description       string `json:"description"`
	Enabled           bool   `json:"enabled"`
	RolloutPercentage int    `json:"rollout_percentage"`
}

type UpdateFeatureFlagRequest struct {
	Name              *string `json:"name"`
	Description       *string `json:"description"`
	Enabled           *bool   `json:"enabled"`
	RolloutPercentage *int    `json:"rollout_percentage"`
}

type PublicFlagResponse struct {
	Enabled bool `json:"enabled"`
}
