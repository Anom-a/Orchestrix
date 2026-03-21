package services

import (
	"errors"
	"os"
	"time"

	"github.com/Anom-a/Orchestrix/internal/database"
	"github.com/Anom-a/Orchestrix/internal/models"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Register(fullName, email, password string) error{
	hashed_password, _ := bcrypt.GenerateFromPassword([]byte(password),10)
	user := models.User{
		FullName: fullName,
		Email: email,
		Password: string(hashed_password),
		Created_at: time.Now(),
	}
	return database.DB.Create(&user).Error
}

func Login(email, password string)(string, error){
	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	if err != nil{
		return "", errors.New("invalid credentails")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil{
		return "", errors.New("invalid incredentails")
	}
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}


