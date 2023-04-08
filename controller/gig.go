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

func CreateGig(c echo.Context) (err error) {
	gig := new(model.Gig)
	if err = c.Bind(gig); err != nil {
		return
	}

	gig.HirerID = GetIDFromToken(c)
	role := GetRoleFromToken(c)

	// Validation
	if gig.Desc == "" || gig.WorkerID == 0 {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Invalid description or id fields"}
	}
	if role != "hirer" {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "Unauthorized gig creation"}
	}

	var user model.User
	result := db.Where("ID = ?", gig.WorkerID).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Worker does not exists"}
	}

	// Scam busting
	gig.Completed = false
	gig.Approved = false

	db.Create(&gig)

	return c.JSON(http.StatusCreated, gig)
}

func GetGig(c echo.Context) (err error) {
	u := new(model.Gig)
	if err = c.Bind(u); err != nil {
		return
	}
	// TODO: ...
	return
}

func UpdateGig(c echo.Context) (err error) {
	u := new(model.Gig)
	if err = c.Bind(u); err != nil {
		return
	}
	// TODO: ...
	return
}

func DeleteGig(c echo.Context) (err error) {

	gigID, _ := strconv.Atoi(c.Param("id"))

	db := database.Connect()
	defer database.Disconnect(db)

	var gig model.Gig
	result := db.Where("ID = ?", gigID).First(&gig)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &echo.HTTPError{Code: http.StatusNotFound, Message: "Gig Not Found"}
	}

	userID := GetIDFromToken(c)
	role := GetRoleFromToken(c)

	if gig.HirerID != userID && role != "admin" {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "Unauthorized delete"}
	}

	db.Delete(&gig)

	return c.NoContent(http.StatusNoContent)
}

func GetPendingGigs(c echo.Context) (err error) {
	userID := GetIDFromToken(c)
	role := GetRoleFromToken(c)

	var gigs []model.Gig

	if role == "worker" {
		db.Where("worker_id = ?", userID).Find(&gigs)
	} else if role == "hirer" {
		db.Where("hirer_id = ?", userID).Find(&gigs)
	} else {
		return &echo.HTTPError{Code: http.StatusTeapot, Message: "I'm a Teapot"}
	}

	return c.JSON(http.StatusOK, gigs)
}

func ApproveGig(c echo.Context) (err error) {
	gigID, _ := strconv.Atoi(c.Param("id"))

	db := database.Connect()
	defer database.Disconnect(db)

	var gig model.Gig
	result := db.Where("ID = ?", gigID).First(&gig)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &echo.HTTPError{Code: http.StatusNotFound, Message: "Gig Not Found"}
	}

	userID := GetIDFromToken(c)
	role := GetRoleFromToken(c)

	if role == "hirer" {
		return &echo.HTTPError{Code: http.StatusForbidden, Message: "If the working class produces everything, then everything belongs to them."}
	}

	if gig.WorkerID != userID {
		return &echo.HTTPError{Code: http.StatusForbidden, Message: "This gig doesn't belong to you."}
	}

	gig.Approved = true

	db.Save(&gig)

	return c.NoContent(http.StatusOK)
}

func CompleteGig(c echo.Context) (err error) {
	gigID, _ := strconv.Atoi(c.Param("id"))

	db := database.Connect()
	defer database.Disconnect(db)

	var gig model.Gig
	result := db.Where("ID = ?", gigID).First(&gig)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &echo.HTTPError{Code: http.StatusNotFound, Message: "Gig Not Found"}
	}

	userID := GetIDFromToken(c)
	role := GetRoleFromToken(c)

	if role == "hirer" {
		return &echo.HTTPError{Code: http.StatusForbidden, Message: "If the working class produces everything, then everything belongs to them."}
	}

	if gig.WorkerID != userID {
		return &echo.HTTPError{Code: http.StatusForbidden, Message: "This gig doesn't belong to you."}
	}

	gig.Completed = true

	db.Save(&gig)

	return c.NoContent(http.StatusOK)
}
