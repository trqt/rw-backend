package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"readyworker.com/backend/database"
	"readyworker.com/backend/model"
)

type ErrorResponse struct {
	Err string
}

type error interface {
	Error() string
}

var db = database.Connect()

func SignUp(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error on parsing request", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Validation
	cpf := r.FormValue("cpf")
	name := r.FormValue("name")
	email := r.FormValue("email")
	desc := r.FormValue("description")

	err = validateCpf(cpf)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	err = validateName(name)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	err = validateEmail(email)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	err = validateDesc(desc)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}

	// Check if the user already exists
	var trash model.User
	result := db.Where("cpf = ?", cpf).Find(&trash)

	if result.RowsAffected > 0 {
		http.Error(w, "User already registered", http.StatusInternalServerError)
		return
	}

	password_hash, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Encryption error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	user := &model.User{
		Type:        "worker",
		Cpf:         cpf,
		Name:        name,
		Password:    password_hash,
		Email:       email,
		Description: desc,
	}

	db.Create(user)

	// Redirect to login
	http.Redirect(w, r, "/login", http.StatusFound) // maybe /login?account_created=1
}

func Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error on parsing request", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	cpf := r.FormValue("cpf")           // Sanitize?
	password := r.FormValue("password") // Sanitize?
	log.Println("User login attempt by: ", cpf)

	// Accept CPF, Email or Username
	var user model.User
	result := db.Where("cpf = ?", cpf).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(w, "User not Found", http.StatusAccepted)
		//http.Redirect(w, r, "/login?retry=1", http.StatusPermanentRedirect)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		http.Error(w, "Hash don't match", http.StatusAccepted)
		//http.Redirect(w, r, "/login?retry=1", http.StatusPermanentRedirect)
		return
	}

	// Check if the user exists or the password is wrong
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(w, "Wrong password", http.StatusAccepted)
	}

	log.Println(time.Now().Format(time.RFC850), "Successful login by: ", cpf)

	w.Write([]byte("Hello " + user.Email))
}

func validateCpf(cpf string) error {
	// Remove all non-digit chars.
	cpf = regexp.MustCompile(`\D`).ReplaceAllString(cpf, "")
	// TODO: Really validate CPF
	if len(cpf) != 11 {
		return errors.New("Invalid CPF")
	}
	return nil
}

func validateName(name string) error {
	// TODO: Sanitise?
	if len(name) > 200 {
		return errors.New("Name too long")
	}
	return nil
}

func validateEmail(email string) error {
	if len(email) > 150 {
		return errors.New("E-mail too long")
	}

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return errors.New("Invalid E-mail")
	}
	return nil
}

func validateDesc(desc string) error {
	// TODO: Sanitise?
	if len(desc) > 500 {
		return errors.New("Description too long")
	}
	return nil
}
