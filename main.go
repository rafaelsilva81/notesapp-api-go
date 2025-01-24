package main

import (
	"fmt"
	"log"
	"net/http"
	"notesapp/api/config"

	"notesapp/api/handlers"
	"os"

	_ "notesapp/api/docs"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	httpSwagger "github.com/swaggo/http-swagger"
)

func startHttpServer() {
	port := os.Getenv("SERVER_PORT")
	databaseUrl := os.Getenv("DATABASE_URL")

	config.InitDatabase(databaseUrl)

	// Quando essa função encerrar, a função db.Close() também será executada para fechar a conexão com o banco de dados.
	defer config.CloseDatabase()

	log.Println("Server started on http://localhost:" + port)
	appRouter := mux.NewRouter()
	appRouter.Use(mux.CORSMethodMiddleware(appRouter))

	// Default middleware (CORS e Json)
	appRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	// Rotas
	handlers.InitializeNotesHandler(appRouter)

	// Swagger
	appRouter.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:"+port+"/swagger/doc.json"),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DomID("swagger-ui"),
	)).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+port, appRouter))

}

// @title Notesapp API docs
// @version 1.0
// @description Documentation for the notesapp API.

// @host
// @BasePath /
func main() {
	fmt.Println("Starting server...")
	startHttpServer()
}
