package app

import (
	"github.com/task-manager/cache"
	"github.com/task-manager/config"
	"github.com/task-manager/db"
)

type App struct {
	config     *config.Config
	postgresDB *db.DB
	redisDB    *cache.Rdb
}

func (app *App) Conf() *config.Config {
	return app.config
}

func (app *App) PostgresDB() *db.DB {
	return app.postgresDB
}

func (app *App) RedisDB() *cache.Rdb {
	return app.redisDB
}

func BuildApp(cfg *config.Config,
	postgres *db.DB,
	redis *cache.Rdb) *App {
	return &App{config: cfg,
		postgresDB: postgres,
		redisDB:    redis}
}
