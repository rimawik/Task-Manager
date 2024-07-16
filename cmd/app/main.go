package main

import (
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/task-manager/app"
	"github.com/task-manager/cache"
	"github.com/task-manager/config"
	"github.com/task-manager/db"
	"github.com/task-manager/routes"
	
	
)

// @title Task API
// @version 1.0
// @description This is a sample server for managing tasks.
// @host localhost:8080
// @BasePath /v1/
func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("couldn't load configuration: %v", err)
	}
	postgresDB, err := db.InitDB(*cfg)
	if err != nil {
		logrus.Fatalf("couldn't initialize db: %v", err)
	}
	redisDB, err := cache.ConnectToRedis(*cfg)
	if err != nil {
		logrus.Fatalf("couldn't connect to redis: %v", err)
	}

	app := app.BuildApp(cfg, postgresDB, redisDB)

	r := routes.NewRouter(app)



	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("could not start server: %v", err)
	}

}
