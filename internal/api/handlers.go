package api

import (
	"errors"
	"net/http"

	"github.com/edteam/feature-flag-service/internal/cache"
	"github.com/edteam/feature-flag-service/internal/db"
	"github.com/edteam/feature-flag-service/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type Handler struct {
	queries *db.Queries
	cache   *cache.FlagCache
}

func NewHandler(queries *db.Queries, flagCache *cache.FlagCache) *Handler {
	return &Handler{queries: queries, cache: flagCache}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", h.Health)

		apps := v1.Group("/apps")
		{
			apps.POST("", h.CreateApplication)
			apps.GET("", h.ListApplications)
			apps.GET("/:appId", h.GetApplication)
			apps.PUT("/:appId", h.UpdateApplication)
			apps.DELETE("/:appId", h.DeleteApplication)

			apps.POST("/:appId/flags", h.CreateFeatureFlag)
			apps.GET("/:appId/flags", h.ListFeatureFlagsByApp)
		}

		v1.GET("/flags/:flagId", h.GetFeatureFlag)
		v1.PUT("/flags/:flagId", h.UpdateFeatureFlag)
		v1.DELETE("/flags/:flagId", h.DeleteFeatureFlag)

		v1.GET("/public/flag/:appId/:flagName", h.GetPublicFlag)
	}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) CreateApplication(c *gin.Context) {
	var req models.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New()
	app, err := h.queries.CreateApplication(c.Request.Context(), db.CreateApplicationParams{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create application"})
		return
	}

	c.JSON(http.StatusCreated, db.ToApplication(app))
}

func (h *Handler) ListApplications(c *gin.Context) {
	apps, err := h.queries.ListApplications(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list applications"})
		return
	}

	result := make([]models.Application, len(apps))
	for i, a := range apps {
		result[i] = db.ToApplication(a)
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetApplication(c *gin.Context) {
	appID, err := parseAppID(c)
	if err != nil {
		return
	}

	app, err := h.queries.GetApplication(c.Request.Context(), appID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}
	c.JSON(http.StatusOK, db.ToApplication(app))
}

func (h *Handler) UpdateApplication(c *gin.Context) {
	appID, err := parseAppID(c)
	if err != nil {
		return
	}

	var req models.UpdateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app, err := h.queries.UpdateApplication(c.Request.Context(), db.UpdateApplicationParams{
		ID:          appID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}
	c.JSON(http.StatusOK, db.ToApplication(app))
}

func (h *Handler) DeleteApplication(c *gin.Context) {
	appID, err := parseAppID(c)
	if err != nil {
		return
	}

	if err := h.queries.DeleteApplication(c.Request.Context(), appID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete application"})
		return
	}
	h.cache.InvalidateApp(appID)
	c.Status(http.StatusNoContent)
}

func (h *Handler) CreateFeatureFlag(c *gin.Context) {
	appID, err := parseAppID(c)
	if err != nil {
		return
	}

	var req models.CreateFeatureFlagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.RolloutPercentage < 0 || req.RolloutPercentage > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rollout_percentage must be between 0 and 100"})
		return
	}

	if _, err := h.queries.GetApplication(c.Request.Context(), appID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	id := uuid.New()
	flag, err := h.queries.CreateFeatureFlag(c.Request.Context(), db.CreateFeatureFlagParams{
		ID:                 id,
		AppID:              appID,
		Name:               req.Name,
		Description:        req.Description,
		Enabled:            req.Enabled,
		RolloutPercentage:  int32(req.RolloutPercentage),
	})
	if err != nil {
		if isUniqueViolation(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "flag name already exists for this application"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create feature flag"})
		return
	}

	model := db.ToFeatureFlag(flag)
	h.cache.Set(appID, &model)
	c.JSON(http.StatusCreated, model)
}

func (h *Handler) ListFeatureFlagsByApp(c *gin.Context) {
	appID, err := parseAppID(c)
	if err != nil {
		return
	}

	flags, err := h.queries.ListFeatureFlagsByApp(c.Request.Context(), appID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list feature flags"})
		return
	}

	result := make([]models.FeatureFlag, len(flags))
	for i, f := range flags {
		result[i] = db.ToFeatureFlag(f)
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetFeatureFlag(c *gin.Context) {
	flagID, err := parseFlagID(c)
	if err != nil {
		return
	}

	flag, err := h.queries.GetFeatureFlag(c.Request.Context(), flagID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "feature flag not found"})
		return
	}
	c.JSON(http.StatusOK, db.ToFeatureFlag(flag))
}

func (h *Handler) UpdateFeatureFlag(c *gin.Context) {
	flagID, err := parseFlagID(c)
	if err != nil {
		return
	}

	var req models.UpdateFeatureFlagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing, err := h.queries.GetFeatureFlag(c.Request.Context(), flagID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "feature flag not found"})
		return
	}

	name := existing.Name
	description := existing.Description
	enabled := existing.Enabled
	rollout := existing.RolloutPercentage

	if req.Name != nil {
		name = *req.Name
	}
	if req.Description != nil {
		description = *req.Description
	}
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	if req.RolloutPercentage != nil {
		if *req.RolloutPercentage < 0 || *req.RolloutPercentage > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "rollout_percentage must be between 0 and 100"})
			return
		}
		rollout = int32(*req.RolloutPercentage)
	}

	flag, err := h.queries.UpdateFeatureFlag(c.Request.Context(), db.UpdateFeatureFlagParams{
		ID:                 flagID,
		Name:               name,
		Description:        description,
		Enabled:            enabled,
		RolloutPercentage:  rollout,
	})
	if err != nil {
		if isUniqueViolation(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "flag name already exists for this application"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update feature flag"})
		return
	}

	model := db.ToFeatureFlag(flag)
	h.cache.Set(model.AppID, &model)
	c.JSON(http.StatusOK, model)
}

func (h *Handler) DeleteFeatureFlag(c *gin.Context) {
	flagID, err := parseFlagID(c)
	if err != nil {
		return
	}

	existing, err := h.queries.GetFeatureFlag(c.Request.Context(), flagID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "feature flag not found"})
		return
	}

	if err := h.queries.DeleteFeatureFlag(c.Request.Context(), flagID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete feature flag"})
		return
	}

	h.cache.Delete(existing.AppID, existing.Name)
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetPublicFlag(c *gin.Context) {
	appID, err := uuid.Parse(c.Param("appId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid app id"})
		return
	}
	flagName := c.Param("flagName")
	userID := c.GetHeader("X-User-Id")

	flag, ok := h.cache.Get(appID, flagName)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "flag not found"})
		return
	}

	enabled := false
	if flag.Enabled {
		enabled = evaluateRollout(flag.Name, userID, flag.RolloutPercentage)
	}

	c.JSON(http.StatusOK, models.PublicFlagResponse{Enabled: enabled})
}

func parseAppID(c *gin.Context) (uuid.UUID, error) {
	appID, err := uuid.Parse(c.Param("appId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid app id"})
	}
	return appID, err
}

func parseFlagID(c *gin.Context) (uuid.UUID, error) {
	flagID, err := uuid.Parse(c.Param("flagId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid flag id"})
	}
	return flagID, err
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
