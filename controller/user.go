package controller

import (
	"errors"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"readyworker.com/backend/database"
	"readyworker.com/backend/model"
)

var db = database.Connect()

func SignUp(c echo.Context) error {
	u := new(model.User)

	if err := c.Bind(u); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "Internal server error"}
	}

	// Validation
	if u.Role == "" || (u.Role != "worker" && u.Role != "hirer") {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Invalid Role"}
	}
	if err := validateCpf(u.Cpf); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if err := validateName(u.Name); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if err := validateEmail(u.Email); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if err := validatePassword(u.Password); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if err := validateDesc(u.Description); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	// Check if the user already exists
	var trash model.User
	result := db.Where("cpf = ?", u.Cpf).Find(&trash)

	if result.RowsAffected > 0 {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "User already registered"}
	}

	password_hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "Internal server error"}
	}

	u.Password = string(password_hash)

	db.Create(u)

	u.Password = ""
	return c.JSON(http.StatusOK, u)
}

func Login(c echo.Context) error {
	u := new(model.User)

	if err := c.Bind(u); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "Internal server error"}
	}

	if u.Email == "" && u.Cpf == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Blank CPF or email"}
	}

	// Accept CPF or E-mail
	var user model.User

	if u.Email == "" {
		result := db.Where("cpf = ?", u.Cpf).First(&user)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: "User not Found"}
		}
	} else if u.Cpf == "" {
		result := db.Where("email = ?", u.Email).First(&user)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: "User not Found"}
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "Wrong CPF/Email or password"}
	}

	// JWT Token generation

	token, err := GenerateJWT(user.ID, user.Role)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"userid": user.ID,
		"token":  token,
	})
}

func GenerateJWT(id uint, role string) (string, error) {
	secret := os.Getenv("JWT_SECRETKEY")

	key := []byte(secret)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["id"] = id
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 45).Unix()

	tokenString, err := token.SignedString(key)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func GetIDFromToken(c echo.Context) uint {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return uint(claims["id"].(float64))
}

func GetRoleFromToken(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["role"].(string)
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
	if name == "" || len(name) < 5 {
		return errors.New("Name too short")
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

func validatePassword(pass string) error {
	// TODO: Sanitise?
	// TODO: Require numbers and special chars
	if len(pass) > 50 {
		return errors.New("Password too long")
	}
	if len(pass) < 8 {
		return errors.New("Password too short")
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
