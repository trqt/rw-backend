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

func CreateComment(c echo.Context) (err error) {
	comment := new(model.Comment)
	if err = c.Bind(comment); err != nil {
		return
	}

	comment.AuthorID = GetIDFromToken(c)
	role := GetRoleFromToken(c)

	// Validation
	if comment.Content == "" || comment.WorkerID == 0 {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Invalid content or id fields"}
	}
	if role != "hirer" {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "Unauthorized comment"}
	}

	// TODO: Check if the hirer hired the worker

	db := database.Connect()
	defer database.Disconnect(db)

	var author model.User
	result := db.Where("ID = ?", comment.AuthorID).First(&author)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Author does not exists"}
	}
	result = db.Where("ID = ?", comment.WorkerID).First(&author)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Worker does not exists"}
	}

	db.Create(&comment)

	return c.JSON(http.StatusCreated, comment)

}
func DeleteComment(c echo.Context) error {
	commentID, _ := strconv.Atoi(c.Param("id"))

	db := database.Connect()
	defer database.Disconnect(db)

	var comment model.Comment
	result := db.Where("ID = ?", commentID).First(&comment)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &echo.HTTPError{Code: http.StatusNotFound, Message: "Comment Not Found"}
	}

	userID := GetIDFromToken(c)
	role := GetRoleFromToken(c)

	if comment.AuthorID != userID && role != "admin" {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "Unauthorized delete"}
	}

	db.Delete(&comment)

	return c.NoContent(http.StatusNoContent)
}
func GetComment(c echo.Context) error {
	commentID, _ := strconv.Atoi(c.Param("id"))

	db := database.Connect()
	defer database.Disconnect(db)

	var comment model.Comment
	result := db.Where("ID = ?", commentID).First(&comment)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &echo.HTTPError{Code: http.StatusNotFound, Message: "Comment Not Found"}
	}

	return c.JSON(http.StatusOK, comment)

}
