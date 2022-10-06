package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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
	u := new(model.Gig)
	if err = c.Bind(u); err != nil {
		return
	}
	// TODO: ...
	return
}

func GetUnapprovedGigs(c echo.Context) (err error) {
	userID := GetIDFromToken(c)
	role := GetRoleFromToken(c)

	var gigs []model.Gig

	if role == "worker" {
		db.Where("worker_id = ? AND approved = ?", userID, false).Find(&gigs)
	} else if role == "hirer" {
		db.Where("hirer_id = ? AND approved = ?", userID, false).Find(&gigs)
	} else {
		return &echo.HTTPError{Code: http.StatusTeapot, Message: "I'm a Teapot"}
	}

	return c.JSON(http.StatusOK, gigs)
}
