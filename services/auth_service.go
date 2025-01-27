package services

import (
	"database/sql"
	"fmt"
	"notesapp/api/config"
	"notesapp/api/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	DB *sql.DB
}

// Struct for JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MaxCost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func generateJWTToken(username string) (string, error) {
	expirationTime := time.Now().Add(config.TokenExpiration)

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JWTSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func NewAuthService() (*AuthService, error) {
	db := config.GetDatabase()
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS users (
				id INTEGER PRIMARY KEY,
				username TEXT NOT NULL UNIQUE,
				password TEXT NOT NULL
    )`

	if _, err := db.Exec(createTableQuery); err != nil {
		return nil, err
	}

	return &AuthService{DB: db}, nil
}

func (service *AuthService) Register(username string, password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	stmt, err := service.DB.Prepare("INSERT INTO users (id, username, password) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(nil, username, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (service *AuthService) Login(username string, password string) (string, error) {
	stmt, err := service.DB.Prepare("SELECT id, username, password FROM users WHERE username = ?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var user models.User
	err = stmt.QueryRow(username).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	// Token JWT
	token, err := generateJWTToken(username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (service *AuthService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return config.JWTSecret, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("Invalid token")
	}

	return claims.Username, nil
}
