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
	if comment.Content == "" || comment.WorkerID == 0 || comment.Rating > 10 || comment.Rating < 0 {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Invalid content, id or rating fields"}
	}
	if role != "hirer" {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "Unauthorized comment"}
	}

	// TODO: Check if the hirer hired the worker

	db := database.Connect()
	defer database.Disconnect(db)

	var user model.User
	/*result := db.Where("ID = ?", comment.AuthorID).First(&author)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Author does not exists"}
	}*/
	result := db.Where("ID = ?", comment.WorkerID).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Worker does not exists"}
	}

	if user.Role != "worker" {
		return &echo.HTTPError{Code: http.StatusForbidden, Message: "worker_id doesn't belongs to an worker"}
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

func GetComments(c echo.Context) error {
	userID, _ := strconv.Atoi(c.Param("id"))

	db := database.Connect()
	defer database.Disconnect(db)

	var comments []model.Comment
	result := db.Where("worker_id = ?", userID).Find(&comments)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &echo.HTTPError{Code: http.StatusNotFound, Message: "User doesn't have comments"}
	}

	return c.JSON(http.StatusOK, comments)

}

func GetRating(c echo.Context) error {
	userID, _ := strconv.Atoi(c.Param("id"))

	db := database.Connect()
	defer database.Disconnect(db)

	var comments []model.Comment
	result := db.Where("worker_id = ?", userID).Find(&comments)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &echo.HTTPError{Code: http.StatusNotFound, Message: "User doesn't have comments"}
	}

	median_rating := 0
	for _, comment := range comments {
		median_rating += int(comment.Rating)
	}
	if len(comments) == 0 {
		return c.JSON(http.StatusBadRequest, "Não há avaliações.")
	}
	median_rating = median_rating / len(comments)

	return c.JSON(http.StatusOK, median_rating)
}
