package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"readyworker.com/backend/database"
	"readyworker.com/backend/model"
)

func GetUsersFromCategory(c echo.Context) error {

	category := c.Param("name")

	db := database.Connect()
	defer database.Disconnect(db)

	var users []model.User
	result := db.Where("category = ?", category).Find(&users)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &echo.HTTPError{Code: http.StatusNotFound, Message: "User Not Found"}
	}

	return c.JSON(http.StatusOK, users)
}
