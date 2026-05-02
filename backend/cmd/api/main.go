package main

import (
	"log"

	"beautymed/internal/config"
	"beautymed/internal/db"
	"beautymed/internal/router"
	"beautymed/internal/service"

	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	database := db.Connect(cfg.DatabaseURL)
	defer database.Close()

	db.Migrate(cfg.DatabaseURL)

	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatalf("invalid redis url: %v", err)
	}
	rdb := redis.NewClient(opt)

	svc := service.New(database, rdb, cfg)
	r := router.New(svc, cfg)

	log.Printf("starting server on :%s", cfg.Port)
	log.Fatal(r.Run(":" + cfg.Port))
}
