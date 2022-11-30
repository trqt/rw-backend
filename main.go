package main

import (
	"bufio"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"readyworker.com/backend/controller"
	"readyworker.com/backend/database"
)

func main() {
	e := echo.New()

	// Debug Mode
	e.Debug = true

	// Load careers

	f, err := os.Open("careers.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		// do something with a line
		controller.Categories = append(controller.Categories, scanner.Text())
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(middleware.Gzip())
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(os.Getenv("JWT_SECRETKEY")),
		Skipper: func(c echo.Context) bool {
			// Skip authentication for signup and login requests
			if c.Path() == "/api/login" || c.Path() == "/api/signup" || c.Path() == "/api/categories" {
				return true
			}
			return false
		},
	}))

	// Migration
	db := database.Connect()
	database.Migrate(db)

	e.POST("/api/login", controller.Login)

	e.POST("/api/signup", controller.SignUp)

	e.GET("/api/user/:id", controller.UserPage)

	e.POST("/api/comment", controller.CreateComment)
	e.GET("/api/comment/:id", controller.GetComment)
	e.DELETE("/api/comment/:id", controller.DeleteComment)

	e.GET("/api/comments/:id", controller.GetComments)

	e.POST("/api/gig", controller.CreateGig)
	e.GET("/api/gig/:id", controller.GetGig)
	e.DELETE("/api/gig/:id", controller.DeleteGig)

	e.GET("/api/gigs", controller.GetPendingGigs)
	e.POST("/api/approve/:id", controller.ApproveGig)

	e.GET("/api/category/:name", controller.GetUsersFromCategory)
	e.GET("/api/categories", controller.GetCategories)

	e.GET("/api/rating/:id", controller.GetRating)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	e.Logger.Fatal(e.Start(":8080"))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pages/404.html")
}
