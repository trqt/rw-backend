package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"readyworker.com/backend/database"
	"readyworker.com/backend/model"
)

func PutComment(w http.ResponseWriter, r *http.Request) {
	var comment model.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)

	if err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	db := database.Connect()
	defer database.Disconnect(db)
	// TODO: Authenticate
	var author model.User
	result := db.Where("ID = ?", comment.AuthorID).First(&author)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(w, "Author does not exists", http.StatusBadRequest)
		//http.Redirect(w, r, "/login?retry=1", http.StatusPermanentRedirect)
		return
	}
	result = db.Where("ID = ?", comment.WorkerID).First(&author)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(w, "Worker does not exists", http.StatusBadRequest)
		//http.Redirect(w, r, "/login?retry=1", http.StatusPermanentRedirect)
		return
	}

	db.Create(&comment)

}
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	// TODO: Authenticate
	commentID := mux.Vars(r)["id"]

	db := database.Connect()
	defer database.Disconnect(db)

	var comment model.Comment
	result := db.Where("ID = ?", commentID).First(&comment)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(w, "Comment Not Found", http.StatusNotFound)
		//http.Redirect(w, r, "/login?retry=1", http.StatusPermanentRedirect)
		return
	}
	db.Delete(&comment)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
	})
}
func GetComment(w http.ResponseWriter, r *http.Request) {
	// TODO: Authenticate
	commentID := mux.Vars(r)["id"]

	db := database.Connect()
	defer database.Disconnect(db)

	var comment model.Comment
	result := db.Where("ID = ?", commentID).First(&comment)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(w, "Comment Not Found", http.StatusNotFound)
		//http.Redirect(w, r, "/login?retry=1", http.StatusPermanentRedirect)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&comment)

}
