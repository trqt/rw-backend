package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"readyworker.com/backend/database"
	"readyworker.com/backend/model"
)

func UserPage(c echo.Context) error {
	userID, _ := strconv.Atoi(c.Param("id"))

	db := database.Connect()
	defer database.Disconnect(db)

	var user model.User
	result := db.Where("ID = ?", userID).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &echo.HTTPError{Code: http.StatusNotFound, Message: "User Not Found"}
	}
	// Do not leak private information
	user.Cpf = ""
	user.Email = ""
	user.Phone = ""
	user.Password = ""

	return c.JSON(http.StatusOK, user)
}

func GetContactInfo(c echo.Context) error {
	userID, _ := strconv.Atoi(c.Param("id"))
	role := GetRoleFromToken(c)

	if role != "hirer" {
		return &echo.HTTPError{Code: http.StatusForbidden, Message: "Forbidden"}
	}

	// TODO: Cross check with gig approval

	db := database.Connect()
	defer database.Disconnect(db)

	var user model.User
	result := db.Where("ID = ?", userID).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &echo.HTTPError{Code: http.StatusNotFound, Message: "User Not Found"}
	}
	// Do not leak private information
	user.Cpf = ""
	user.Password = ""

	return c.JSON(http.StatusOK, user)
}
