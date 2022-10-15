package main

import (
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

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(middleware.Gzip())
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(os.Getenv("JWT_SECRETKEY")),
		Skipper: func(c echo.Context) bool {
			// Skip authentication for signup and login requests
			if c.Path() == "/login" || c.Path() == "/signup" {
				return true
			}
			return false
		},
	}))

	// Migration
	db := database.Connect()
	database.Migrate(db)

	e.Static("/static", "public")

	e.POST("/login", controller.Login)

	e.POST("/signup", controller.SignUp)

	e.GET("/user/:id", controller.UserPage)

	e.POST("/comment", controller.CreateComment)
	e.GET("/comment/:id", controller.GetComment)
	e.DELETE("/comment/:id", controller.DeleteComment)

	e.GET("/comments/:id", controller.GetComments)

	e.POST("/gig", controller.CreateGig)
	e.GET("/gig/:id", controller.GetGig)
	e.DELETE("/gig/:id", controller.DeleteGig)

	e.GET("/gigs", controller.GetUnapprovedGigs)

	e.GET("/category/:name", controller.GetUsersFromCategory)

	e.Logger.Fatal(e.Start(":8080"))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pages/404.html")
}
