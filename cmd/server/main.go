package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/edteam/feature-flag-service/internal/api"
	"github.com/edteam/feature-flag-service/internal/cache"
	"github.com/edteam/feature-flag-service/internal/db"
	"github.com/edteam/feature-flag-service/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	port := getEnv("PORT", "8080")
	databaseURL := getEnv("DATABASE_URL", "postgres://flaguser:flagpass@localhost:5432/feature_flags?sslmode=disable")
	migrationsPath := getEnv("MIGRATIONS_PATH", "file://migrations")
	cacheRefreshSec, _ := strconv.Atoi(getEnv("CACHE_REFRESH_INTERVAL", "30"))
	ginMode := getEnv("GIN_MODE", "debug")

	if os.Getenv("SKIP_MIGRATIONS") != "true" {
		if err := db.RunMigrations(databaseURL, migrationsPath); err != nil {
			log.Fatalf("migrations failed: %v", err)
		}
		log.Println("migrations applied successfully")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := db.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()

	queries := db.New(pool)
	flagCache := cache.NewFlagCache(queries)
	flagCache.StartRefreshLoop(ctx, time.Duration(cacheRefreshSec)*time.Second)

	gin.SetMode(ginMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	handler := api.NewHandler(queries, flagCache)
	handler.RegisterRoutes(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("server listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Println("server stopped")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
