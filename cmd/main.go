package main

import (
	"log"
	"os"

	"short-link-service/internal/config"
	"short-link-service/internal/handler"
	"short-link-service/internal/logger"
	"short-link-service/internal/repo"
	"short-link-service/internal/router"
	"short-link-service/internal/service"
	"short-link-service/internal/store"
)

func main() {
	configPath := "configs/config.example.yaml"
	if v := os.Getenv("APP_CONFIG"); v != "" {
		configPath = v
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("读取配置失败: %v", err)
	}

	logger.Init(cfg.App.Name)

	db, err := store.InitMySQL(cfg.MySQL.DSN)
	if err != nil {
		log.Fatalf("初始化 MySQL 失败: %v", err)
	}

	cache, err := store.InitRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Printf("初始化 Redis 失败，将降级为无缓存模式: %v", err)
		cache = nil
	}

	linkRepo := repo.NewLinkRepo(db, cache)
	linkService := &service.ShortLinkService{Repo: linkRepo}
	h := &handler.Handler{Service: linkService}

	r := router.Setup(h)
	if err := r.Run(cfg.HTTP.Addr); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
