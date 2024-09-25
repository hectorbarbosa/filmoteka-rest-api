package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "filmoteka/docs"
	"filmoteka/internal/app/service"
	"filmoteka/internal/restapi"
	"filmoteka/internal/storage/postgresql"
)

//	@title			Swagger filmoteka API
//	@version		1.0
//	@description	This is a sample filmoteka server.
//  @schemes        http
//  @host           localhost:8080

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db_url := os.Getenv("DB_URL")
	db := newDB(db_url)
	defer db.Close()

	repoFilms := postgresql.NewFilm(db)           // Film Repository
	svcFilms := service.NewFilmService(repoFilms) // Film Service

	repoActors := postgresql.NewActor(db)            // Actor Repository
	svcActors := service.NewActorService(repoActors) // Actor Service

	r := mux.NewRouter()

	restapi.NewFilmHandler(svcFilms).Register(r)
	restapi.NewActorHandler(svcActors).Register(r)

	address := os.Getenv("BIND_ADDR")
	// swagUrl := "http://localhost:" + address + "docs/doc.json"
	swagUrl := "./docs/doc.json"

	r.PathPrefix("/docs/").Handler(httpSwagger.Handler(
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

func newDB(dbUrl string) *sql.DB {

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalln("Couldn't open DB", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalln("Couldn't ping DB", err)
	}

	return db
}
