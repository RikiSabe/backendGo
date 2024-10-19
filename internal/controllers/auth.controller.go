package controllers

import (
	"backend/internal/db"
	"backend/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type auth struct {
}

var secretKey = []byte(os.Getenv("JWT_SECRET"))
var Auth auth

func (auth) AuthLogin(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la solicitud para obtener el username y password
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Solicitud inválida"})
		return
	}

	var userR models.Usuario
	tx := db.GDB.Begin()

	// Buscar el usuario por su nombre de usuario
	if err := tx.Where("usuario = ? and rol = 'lecturador'", user.Username).First(&userR).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Usuario no encontrado"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al recuperar el usuario"})
		return
	}

	tx.Commit()

	// Verificar la contraseña
	err := bcrypt.CompareHashAndPassword([]byte(userR.Contra), []byte(user.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Contraseña incorrecta"})
		return
	}

	// Crear el token JWT
	token, err := createToken(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al crear el token"})
		return
	}

	// Enviar el token como respuesta JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (auth) AuthLoginWeb(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la solicitud para obtener el username y password
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Solicitud inválida"})
		return
	}

	var userR models.Usuario
	tx := db.GDB.Begin()

	// Buscar el usuario por su nombre de usuario
	if err := tx.Where("usuario = ? and rol = 'admin'", user.Username).First(&userR).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Usuario no encontrado"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al recuperar el usuario"})
		return
	}

	tx.Commit()

	// Verificar la contraseña
	err := bcrypt.CompareHashAndPassword([]byte(userR.Contra), []byte(user.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Contraseña incorrecta"})
		return
	}

	// Crear el token JWT
	token, err := createToken(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al crear el token"})
		return
	}

	// Enviar el token como respuesta JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// createToken creates a new JWT token for a given username
func createToken(username string) (string, error) {
	// Create a new JWT token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expiration time
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyToken verifies a given JWT token and returns the parsed token if valid
func verifyToken(tokenString string) (*jwt.Token, error) {
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	// Check for parsing errors
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Return the verified token
	return token, nil
}
