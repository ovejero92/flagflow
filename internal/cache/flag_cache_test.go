package cache

import (
	"testing"

	"github.com/edteam/feature-flag-service/internal/models"
	"github.com/google/uuid"
)

func TestFlagCache_ConcurrentAccess(t *testing.T) {
	c := NewFlagCache(nil)
	appID := uuid.New()
	flag := &models.FeatureFlag{
		ID:                uuid.New(),
		AppID:             appID,
		Name:              "test-flag",
		Enabled:           true,
		RolloutPercentage: 100,
	}

	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			c.Set(appID, flag)
			c.Get(appID, "test-flag")
			c.GetAllForApp(appID)
			done <- struct{}{}
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}

	got, ok := c.Get(appID, "test-flag")
	if !ok || got.Name != "test-flag" {
		t.Fatal("expected flag in cache")
	}

	c.Delete(appID, "test-flag")
	if _, ok := c.Get(appID, "test-flag"); ok {
		t.Fatal("expected flag deleted")
	}

	c.Set(appID, flag)
	c.InvalidateApp(appID)
	if _, ok := c.Get(appID, "test-flag"); ok {
		t.Fatal("expected app invalidated")
	}
}
