package authcontroller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/autumnleaf-ra/gorilla-mux/config"
	"github.com/autumnleaf-ra/gorilla-mux/helper"
	"github.com/autumnleaf-ra/gorilla-mux/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ResponseJson = helper.ResponseJson
var ResponseError = helper.ResponseError

func Login(w http.ResponseWriter, r *http.Request) {
	// mengambil input client
	var userInput models.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		ResponseJson(w, http.StatusBadRequest, response)
		return
	}

	defer r.Body.Close()

	// ambil data username
	var user models.User
	if err := models.DB.Where("username = ?", userInput.Username).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response := map[string]string{"message": "Username atau password salah"}
			ResponseJson(w, http.StatusUnauthorized, response)
			return
		default:
			response := map[string]string{"message": err.Error()}
			ResponseJson(w, http.StatusInternalServerError, response)
			return
		}
	}

	//  password checking
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		response := map[string]string{"message": "Username atau password salah"}
		ResponseJson(w, http.StatusUnauthorized, response)
		return
	}

	// proses pembuatan token jwt
	// expired time
	expTime := time.Now().Add(time.Minute * 1)
	claims := &config.JWTClaim{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "gorilla-mux",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	// deklarasi algoritma untuk login
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign token
	token, err := tokenAlgo.SignedString(config.JWT_KEY)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		ResponseJson(w, http.StatusInternalServerError, response)
		return
	}

	// set token to cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token,
		HttpOnly: true,
	})

	response := map[string]string{"message": "Login berhasil !"}
	ResponseJson(w, http.StatusOK, response)
}

func Register(w http.ResponseWriter, r *http.Request) {
	// mengambil input client
	var userInput models.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		ResponseJson(w, http.StatusBadRequest, response)
		return
	}

	defer r.Body.Close()

	// hash password bcrypt
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	userInput.Password = string(hashPassword)

	// insert database
	if err := models.DB.Create(&userInput).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		ResponseJson(w, http.StatusBadRequest, response)
		return
	}

	ResponseJson(w, http.StatusOK, "Akun berhasil dibuat!")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// set token to cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	response := map[string]string{"message": "Logout berhasil!"}
	ResponseJson(w, http.StatusOK, response)
}
