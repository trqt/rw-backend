package controller

import (
	"encoding/json"
	"net/http"

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

	db.Create(&model.Comment{})

}
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	// TODO: ...
}
func GetComment(w http.ResponseWriter, r *http.Request) {
	// TODO: ...
}
