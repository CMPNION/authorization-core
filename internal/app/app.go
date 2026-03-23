package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cmpnion.space/internal/config"
	"cmpnion.space/internal/delivery/httpdelivery"
	"cmpnion.space/internal/repository/jwt"
	"cmpnion.space/internal/repository/password"
	postgresRepo "cmpnion.space/internal/repository/postgres"
	redisRepo "cmpnion.space/internal/repository/redis"
	"cmpnion.space/internal/usecase/auth"

	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
)

type App struct {
	server *http.Server
	pgPool *pgxpool.Pool
	redis  *goredis.Client
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	pgPool, err := pgxpool.New(ctx, cfg.PostgresDSN())
	if err != nil {
		return nil, fmt.Errorf("postgres connect error: %w", err)
	}

	if err := pgPool.Ping(ctx); err != nil {
		pgPool.Close()
		return nil, fmt.Errorf("postgres ping error: %w", err)
	}

	redisClient := goredis.NewClient(&goredis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		redisClient.Close()
		pgPool.Close()
		return nil, fmt.Errorf("redis ping error: %w", err)
	}

	userRepo := postgresRepo.NewUserRepository(pgPool)
	tokenRepo := redisRepo.NewTokenRepository(redisClient)
	hasher := password.NewBcryptHasher(12)
	jwtManager := jwt.NewManager(cfg.JWTSecret)

	authService := auth.NewService(userRepo, tokenRepo, hasher, jwtManager, cfg.TokenTTL)
	handler := httpdelivery.NewAuthHandler(authService)
	router := httpdelivery.NewRouter(handler, authService)

	server := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return &App{
		server: server,
		pgPool: pgPool,
		redis:  redisClient,
	}, nil
}

func (a *App) Start() error {
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	serverErr := a.server.Shutdown(ctx)
	a.redis.Close()
	a.pgPool.Close()
	return serverErr
}
