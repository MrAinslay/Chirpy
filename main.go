package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	mainRouter := chi.NewRouter()
	apiRouter := chi.NewRouter()
	adminRouter := chi.NewRouter()

	db, _ := NewDB("database.json")

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	mainRouter.Handle("/app/*", fsHandler)
	mainRouter.Handle("/app", fsHandler)
	mainRouter.Mount("/api", apiRouter)
	mainRouter.Mount("/admin", adminRouter)

	apiRouter.Get("/reset", apiCfg.handlerReset)
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/chirps", db.getHandler)
	apiRouter.Post("/users", db.usersPostHandler)
	apiRouter.Post("/chirps", db.postHandler)
	apiRouter.Route("/chirps/{id}", func(r chi.Router) {
		r.Get("/", db.getIdHandler)
	})

	adminRouter.Get("/metrics", apiCfg.handlerMetrics)

	corsMux := middlewareCors(mainRouter)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	t, err := template.ParseFiles("metricsTemplate.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.Execute(w, map[string]interface{}{"Metrics": fmt.Sprint(cfg.fileserverHits)}); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}
