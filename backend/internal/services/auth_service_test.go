package services

import (
	"os"
	"testing"
	"time"

	"github.com/Anom-a/Orchestrix/internal/database"
	"github.com/Anom-a/Orchestrix/internal/models"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAuthTestDB(t *testing.T) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite test db: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("failed to migrate user model: %v", err)
	}

	database.DB = db
}

func TestRegisterCreatesUserWithHashedPassword(t *testing.T) {
	setupAuthTestDB(t)

	err := Register("Ada Lovelace", "ada@example.com", "super-secret")
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	var user models.User
	err = database.DB.Where("email = ?", "ada@example.com").First(&user).Error
	if err != nil {
		t.Fatalf("user was not persisted: %v", err)
	}

	if user.Password == "super-secret" {
		t.Fatalf("password should be hashed, got plain text")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("super-secret")); err != nil {
		t.Fatalf("stored hash does not match source password: %v", err)
	}

	if user.Created_at.IsZero() {
		t.Fatalf("created_at should be set")
	}
}

func TestLoginReturnsSignedJWT(t *testing.T) {
	setupAuthTestDB(t)
	t.Setenv("JWT_SECRET", "test-secret")

	hash, err := bcrypt.GenerateFromPassword([]byte("pw123456"), 10)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	user := models.User{
		FullName:   "Grace Hopper",
		Email:      "grace@example.com",
		Password:   string(hash),
		Created_at: time.Now(),
	}
	if err := database.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	tokenString, err := Login("grace@example.com", "pw123456")
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		t.Fatalf("failed to parse jwt: %v", err)
	}

	if !parsedToken.Valid {
		t.Fatalf("expected token to be valid")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("expected MapClaims")
	}

	if claims["user_id"] != float64(user.ID) {
		t.Fatalf("expected user_id claim %d, got %v", user.ID, claims["user_id"])
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	setupAuthTestDB(t)
	_, err := Login("missing@example.com", "whatever")
	if err == nil {
		t.Fatalf("expected error for unknown user")
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	setupAuthTestDB(t)
	t.Setenv("JWT_SECRET", "test-secret")

	hash, _ := bcrypt.GenerateFromPassword([]byte("right-password"), 10)
	seeded := models.User{
		FullName:   "Linus Torvalds",
		Email:      "linus@example.com",
		Password:   string(hash),
		Created_at: time.Now(),
	}
	if err := database.DB.Create(&seeded).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	_, err := Login("linus@example.com", "wrong-password")
	if err == nil {
		t.Fatalf("expected invalid password error")
	}
}
