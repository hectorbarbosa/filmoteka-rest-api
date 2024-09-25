package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "filmoteka/docs"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter()

	// address := os.Getenv("BIND_ADDR")
	address := "localhost:5000"
	// swagUrl := "http://" + address + "/doc.json"
	swagUrl := "./docs/doc.json"

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(swagUrl), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	srv := &http.Server{
		Handler:           r,
		Addr:              address,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}

	log.Println("Starting server ...", address)

	log.Fatal(srv.ListenAndServe())
}
