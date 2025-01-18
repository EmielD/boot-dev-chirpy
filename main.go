package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/emiel/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func (cfg *ApiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		fmt.Println(cfg.fileserverHits.Load())
		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()

	const (
		port       = ":8080"
		assetsPath = "/app/assets/"
		adminPath  = "/admin/"
		appPath    = "/app/"
		staticDir  = "./assets"
	)

	apiCfg := ApiConfig{}
	mux := http.NewServeMux()

	// set up connection to postgres database
	dbUrl := os.Getenv("DB_URL")
	fmt.Println(dbUrl)
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dBQueries := database.New(db)
	apiCfg.db = dBQueries

	// create the file server handler
	fileServer := http.FileServer(http.Dir(staticDir))
	strippedFileServer := http.StripPrefix(assetsPath, fileServer)

	// create the index handler
	indexHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// add routes with middleware
	mux.Handle("GET "+assetsPath, apiCfg.middlewareMetricsInc(strippedFileServer))
	mux.Handle("GET "+appPath, apiCfg.middlewareMetricsInc(indexHandler))

	mux.HandleFunc("GET "+adminPath, adminTask(&apiCfg))
	mux.HandleFunc("POST "+adminPath+"reset", resetTask(&apiCfg))

	mux.HandleFunc("POST /api/login", loginTask(&apiCfg))
	mux.HandleFunc("POST /api/users", createUserTask(&apiCfg))
	mux.HandleFunc("POST /api/chirps", ChirpsTask(&apiCfg))
	mux.HandleFunc("GET /api/chirps", GetChirpsTask(&apiCfg))
	mux.HandleFunc("GET /api/chirps/{chirpID}", GetChirpTask(&apiCfg))
	mux.HandleFunc("GET /api/healthz", HealthTask)

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	fmt.Println("Starting server..")
	log.Fatal(server.ListenAndServe())
}

func resetTask(cfg *ApiConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		cfg.db.Reset(context.Background())
		cfg.fileserverHits.Store(0)
	}
}

func adminTask(cfg *ApiConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/html")

		template :=
			fmt.Sprintf(`<html>
 				 			<body>
    							<h1>Welcome, Chirpy Admin</h1>
    							<p>Chirpy has been visited %d times!</p>
  							</body>
						</html>`, cfg.fileserverHits.Load())

		w.Write([]byte(template))
	}
}
