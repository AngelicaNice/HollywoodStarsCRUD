package main

import (
	"log"
	"net/http"

	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/repository/psql"
	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/service"
	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/transport/rest"
	"github.com/AngelicaNice/HollywoodStarsCRUD/pkg/database"

	_ "github.com/lib/pq"
)

func main() {
	// init db
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

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
