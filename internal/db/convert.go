package db

import (
	"github.com/edteam/feature-flag-service/internal/models"
	"github.com/google/uuid"
)

func ToApplication(a Application) models.Application {
	return models.Application{
		ID:          a.ID,
		Name:        a.Name,
		Description: a.Description,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

func ToFeatureFlag(f FeatureFlag) models.FeatureFlag {
	return models.FeatureFlag{
		ID:                f.ID,
		AppID:             f.AppID,
		Name:              f.Name,
		Description:       f.Description,
		Enabled:           f.Enabled,
		RolloutPercentage: int(f.RolloutPercentage),
		CreatedAt:         f.CreatedAt,
		UpdatedAt:         f.UpdatedAt,
	}
}

func ToFeatureFlags(flags []FeatureFlag) []*models.FeatureFlag {
	result := make([]*models.FeatureFlag, len(flags))
	for i, f := range flags {
		flag := ToFeatureFlag(f)
		result[i] = &flag
	}
	return result
}

func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
