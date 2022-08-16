package main

import (
	"log"
	"net/http"
	"time"

	// "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"readyworker.com/backend/controller"
	"readyworker.com/backend/database"
)

//var SECRET_KEY = []byte("gosecretkey")

func main() {
	log.Println("Starting...")

	database.Connect()

	r := mux.NewRouter()
	r.HandleFunc("/", controller.Home).Methods("GET")
	r.HandleFunc("/login", controller.ServeLogin).Methods("GET")
	r.HandleFunc("/login", controller.Login).Methods("POST")
	//r.HandleFunc("/list", controller.ListAll).Methods("GET")
	r.HandleFunc("/list/{id}", controller.List).Methods("GET")

	r.HandleFunc("/signup", controller.ServeSignUp).Methods("GET")
	r.HandleFunc("/signup", controller.SignUp).Methods("POST")

	r.Use(RequestLoggerMiddleware(r))

	log.Fatal(http.ListenAndServe(":8080", r))
}

func RequestLoggerMiddleware(r *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()
			defer func() {
				log.Printf(
					"[%s] [%v] %s %s %s",
					req.Method,
					time.Since(start),
					req.Host,
					req.URL.Path,
					req.URL.RawQuery,
				)
			}()

			next.ServeHTTP(w, req)
		})
	}
}
