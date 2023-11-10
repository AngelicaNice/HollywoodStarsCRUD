package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/config"
	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/repository/psql"
	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/service"
	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/transport/rest"
	"github.com/AngelicaNice/HollywoodStarsCRUD/pkg/database"

	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

const (
	CONFIG_DIR  = "configs"
	CONFIG_FILE = "main"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

//	@title			Swagger HollywoodStars App API
//	@version		1.0
//	@description	API server for HollywoodStars Application.

// @host		localhost:8080
// @BasePath	/actors
func main() {
	cfg, err := config.NewConfig(CONFIG_DIR, CONFIG_FILE)
	if err != nil {
		log.WithField("config | env", "wrong config | env").Fatal(err)
	}

	db, err := database.CreateDBConnection(
		database.ConnectionInfo{
			Host:     cfg.DB.Host,
			Port:     cfg.DB.Port,
			Username: cfg.DB.Username,
			DBName:   cfg.DB.Name,
			SSLMode:  cfg.DB.SSLMode,
			Password: cfg.DB.Password,
		})
	if err != nil {
		log.WithField("db connection", "failed").Fatal(err)
	}
	defer db.Close()

	actorsRepo := psql.NewActors(db)
	actorsService := service.NewActors(actorsRepo)
	handler := rest.NewHandler(actorsService)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: handler.InitRouter(),
	}

	log.Info("SERVER STARTED")

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
