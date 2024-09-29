package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	esv7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"

	_ "filmoteka/docs"
	elasticsearch "filmoteka/internal/app/elacticsearch"
	"filmoteka/internal/app/service"
	"filmoteka/internal/envvar"
	"filmoteka/internal/restapi"
	"filmoteka/internal/storage/postgresql"
)

//	@title			Swagger filmoteka API
//	@version		1.0
//	@description	This is a sample filmoteka server.
//  @schemes        http
//  @host           localhost:8080

func main() {
	var env string

	flag.StringVar(&env, "env", "", "Environment Variables filename")
	flag.Parse()

	errC, err := run(env)
	if err != nil {
		log.Fatalf("Couldn't run: %s", err)
	}

	if err := <-errC; err != nil {
		log.Fatalf("Error while running: %s", err)
	}

}

func run(env string) (<-chan error, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("zap.NewProduction %w", err)
	}

	if err := envvar.Load(env); err != nil {
		return nil, fmt.Errorf("envvar.Load %w", err)
	}

	conf := envvar.New()

	db, err := newDB(conf)
	if err != nil {
		return nil, fmt.Errorf("newDB %w", err)
	}

	es, err := newElasticSearch(conf)
	if err != nil {
		return nil, fmt.Errorf("newElasticSearch %w", err)
	}

	logging := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(r.Method,
				zap.Time("time", time.Now()),
				zap.String("url", r.URL.String()),
			)

			h.ServeHTTP(w, r)
		})
	}

	serverAddr, err := conf.Get("BIND_ADDR")
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	errC := make(chan error, 1)

	srv := newServer(conf, db, es, logging)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-ctx.Done()

		logger.Info("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer func() {
			logger.Sync()
			db.Close()
			stop()
			cancel()
			close(errC)
		}()

		srv.SetKeepAlivesEnabled(false)

		if err := srv.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}

		logger.Info("Shutdown completed")
	}()

	go func() {
		logger.Info("Listening and serving", zap.String("address", serverAddr))

		// "ListenAndServe always returns a non-nil error. After Shutdown or Close, the returned error is
		// ErrServerClosed."
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errC <- err
		}
	}()

	return errC, nil
}

func newServer(conf *envvar.Configuration, db *sql.DB, es *esv7.Client, mws ...mux.MiddlewareFunc) *http.Server {
	r := mux.NewRouter()

	for _, mw := range mws {
		r.Use(mw)
	}

	repoFilms := postgresql.NewFilm(db) // Film Repository
	repoSearch := elasticsearch.NewFilmSearchRepo(es)
	svcFilms := service.NewFilmService(repoFilms, repoSearch) // Film Service

	repoActors := postgresql.NewActor(db)            // Actor Repository
	svcActors := service.NewActorService(repoActors) // Actor Service

	restapi.NewFilmHandler(svcFilms).Register(r)
	restapi.NewActorHandler(svcActors).Register(r)

	address, err := conf.Get("BIND_ADDR")
	if err != nil {
		log.Fatal(err.Error())
	}

	swagUrl, err := conf.Get("SWAG_URL")
	if err != nil {
		log.Fatal(err.Error())
	}

	// swagUrl := "./docs/doc.json"

	r.PathPrefix("/docs/").Handler(httpSwagger.Handler(
		httpSwagger.URL(swagUrl), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	return &http.Server{
		Handler:           r,
		Addr:              address,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}
}

func newDB(conf *envvar.Configuration) (*sql.DB, error) {
	dbUrl, err := conf.Get("DB_URL")
	if err != nil {
		log.Fatal("No Database url. Check config file")
		return nil, fmt.Errorf("No Database url %w", err)
	}
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("sql.Open %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping %w", err)
	}

	return db, nil
}

func newElasticSearch(conf *envvar.Configuration) (*esv7.Client, error) {
	es, err := esv7.NewDefaultClient()
	if err != nil {
		return nil, fmt.Errorf("elasticsearch.Open %w", err)
	}

	res, err := es.Info()
	if err != nil {
		return nil, fmt.Errorf("es.Info %w", err)
	}
	defer res.Body.Close()

	return es, nil
}
