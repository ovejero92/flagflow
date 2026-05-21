package cache

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/edteam/feature-flag-service/internal/db"
	"github.com/edteam/feature-flag-service/internal/models"
	"github.com/google/uuid"
)

type FlagCache struct {
	mu    sync.RWMutex
	store map[uuid.UUID]map[string]*models.FeatureFlag
	queries *db.Queries
}

func NewFlagCache(queries *db.Queries) *FlagCache {
	return &FlagCache{
		store:   make(map[uuid.UUID]map[string]*models.FeatureFlag),
		queries: queries,
	}
}

func (c *FlagCache) Get(appID uuid.UUID, flagName string) (*models.FeatureFlag, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	appFlags, ok := c.store[appID]
	if !ok {
		return nil, false
	}
	flag, ok := appFlags[flagName]
	return flag, ok
}

func (c *FlagCache) GetAllForApp(appID uuid.UUID) ([]*models.FeatureFlag, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	appFlags, ok := c.store[appID]
	if !ok {
		return nil, false
	}
	result := make([]*models.FeatureFlag, 0, len(appFlags))
	for _, flag := range appFlags {
		result = append(result, flag)
	}
	return result, true
}

func (c *FlagCache) Set(appID uuid.UUID, flag *models.FeatureFlag) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.store[appID]; !ok {
		c.store[appID] = make(map[string]*models.FeatureFlag)
	}
	c.store[appID][flag.Name] = flag
}

func (c *FlagCache) Delete(appID uuid.UUID, flagName string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if appFlags, ok := c.store[appID]; ok {
		delete(appFlags, flagName)
	}
}

func (c *FlagCache) InvalidateApp(appID uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, appID)
}

func (c *FlagCache) Refresh(ctx context.Context) error {
	flags, err := c.queries.ListAllFeatureFlags(ctx)
	if err != nil {
		return err
	}

	newStore := make(map[uuid.UUID]map[string]*models.FeatureFlag)
	for _, f := range flags {
		flag := db.ToFeatureFlag(f)
		if _, ok := newStore[flag.AppID]; !ok {
			newStore[flag.AppID] = make(map[string]*models.FeatureFlag)
		}
		newStore[flag.AppID][flag.Name] = &flag
	}

	c.mu.Lock()
	c.store = newStore
	c.mu.Unlock()
	return nil
}

func (c *FlagCache) StartRefreshLoop(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		if err := c.Refresh(ctx); err != nil {
			log.Printf("cache: initial refresh failed: %v", err)
		}

		for {
			select {
			case <-ctx.Done():
				log.Println("cache: refresh loop stopped")
				return
			case <-ticker.C:
				if err := c.Refresh(ctx); err != nil {
					log.Printf("cache: refresh failed: %v", err)
				}
			}
		}
	}()
}
