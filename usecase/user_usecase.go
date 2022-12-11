package usecase

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/letenk/pokedex/models/web"
	"github.com/letenk/pokedex/repository"
	"github.com/letenk/pokedex/util"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Login(ctx context.Context, req web.UserLoginRequest) (string, error)
}

type Claim struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type userUsecase struct {
	repository repository.UserRepository
}

func NewUsecaseUser(repository repository.UserRepository) *userUsecase {
	return &userUsecase{repository}
}

func (s *userUsecase) Login(ctx context.Context, req web.UserLoginRequest) (string, error) {
	// Get payload
	username := req.Username
	password := req.Password

	// Find user by username
	user, err := s.repository.FindByUsername(ctx, username)
	if user.ID == "" {
		return "", errors.New("username or password incorrect")
	}

	// Other error
	if err != nil {
		log.Println(err)
	}

	// If user is available, compare password hash with password from request use bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("username or password incorrect")
	}

	// If username and password is matched, generate token
	// Create 1 day
	expirationTime := time.Now().AddDate(0, 0, 1)

	// Create clain for payload token
	claim := Claim{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	// Load Config
	config, err := util.LoadConfig("../.")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Signed token with secret key
	signedToken, err := token.SignedString([]byte(config.JWT_SECRET_KEY))
	if err != nil {
		return signedToken, err
	}

	// If success, return token
	return signedToken, nil
}
