package main

import (
	"net/http"
	"os"

	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/repository/psql"
	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/service"
	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/transport/rest"
	"github.com/AngelicaNice/HollywoodStarsCRUD/pkg/database"

	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
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
	db, err := database.CreateDBConnection(
		database.ConnectionInfo{
			Host:     "0.0.0.0",
			Port:     5432,
			Username: "postgres",
			DBName:   "postgres",
			SSLMode:  "disable",
			Password: "goLANGni1nja",
		})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	actorsRepo := psql.NewActors(db)
	actorsService := service.NewActors(actorsRepo)
	handler := rest.NewHandler(actorsService)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler.InitRouter(),
	}

	log.Info("SERVER STARTED")

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
