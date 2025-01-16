package main

import (
	"biggest-mars-pictures/internal/app/clients/nasa"
	"biggest-mars-pictures/internal/app/config"
	"biggest-mars-pictures/internal/app/repository"
	"biggest-mars-pictures/internal/app/services"
	"biggest-mars-pictures/internal/app/transport/httpserver"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if os.Getenv("DEBUG") != "" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	cfg, err := config.Read()
	if err != nil {
		return fmt.Errorf("failde to read config: %w", err)
	}

	client := nasa.NewClient(cfg.NasaAPIKey)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=postgres password=mysecretpassword port=5432 dbname=postgres sslmode=disable",
		PreferSimpleProtocol: true,
	}))
	if err != nil {
		return fmt.Errorf("failde to connect to NASA API: %w", err)
	}
	err = db.AutoMigrate(&repository.Image{})
	if err != nil {
		return err
	}
	lps := services.NewLargestPictureService(client, repository.NewImageRepository(db))

	httpServer := httpserver.NewHttpServer(lps)
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Largest Picture Service API v0.1"))
	}).Methods("GET")
	router.HandleFunc("/lps", httpServer.GetLargestPicture)

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("ListenAndServe: ", err)
	}
	return nil
}
