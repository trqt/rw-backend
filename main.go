package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"readyworker.com/backend/controller"
	"readyworker.com/backend/database"
)

func main() {
	log.Println("Starting...")

	// Migration
	db := database.Connect()
	database.Migrate(db)

	r := mux.NewRouter()
	r.HandleFunc("/", controller.Home).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./pages/static/"))))

	r.HandleFunc("/login", controller.ServeLogin).Methods("GET")
	r.HandleFunc("/login", controller.Login).Methods("POST")
	//r.HandleFunc("/list", controller.ListAll).Methods("GET")
	r.HandleFunc("/list/{id}", controller.List).Methods("GET")

	r.HandleFunc("/signup", controller.ServeSignUp).Methods("GET")
	r.HandleFunc("/signup", controller.SignUp).Methods("POST")

	r.HandleFunc("/comment", controller.PutComment).Methods("PUT")
	r.HandleFunc("/comment", controller.GetComment).Methods("GET")
	r.HandleFunc("/comment", controller.DeleteComment).Methods("DELETE")

	r.Use(RequestLoggerMiddleware(r))

	log.Fatal(http.ListenAndServe(":8080", r))
}

func RequestLoggerMiddleware(r *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()
			defer func() {
				log.Printf(
					"[%s] [%v] %s %s on %s %s %s",
					req.Method,
					time.Since(start),
					req.RemoteAddr,
					req.UserAgent(),
					req.Host,
					req.URL.Path,
					req.URL.RawQuery,
				)
			}()

			next.ServeHTTP(w, req)
		})
	}
}
