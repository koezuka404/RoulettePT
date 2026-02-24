package usecase

import (
	"os"
	"time"

	"backend/model"
	"backend/repository"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	SignUp(user model.User) (model.ResponseUser, error)
	Login(user model.User) (string, error)
}
type userUsecase struct {
	ur repository.IUserRepository
}

func NewUserUsecase(ur repository.IUserRepository) IUserUsecase {
	return &userUsecase{ur}
}

func (uu *userUsecase) SignUp(user model.User) (model.ResponseUser, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.ResponseUser{}, err
	}
	newUser := model.User{Email: user.Email, Password: string(hash)}
	if err := uu.ur.Create(&newUser); err != nil {
		return model.ResponseUser{}, err
	}
	resUser := model.ResponseUser{
		ID:    newUser.ID,
		Email: newUser.Email,
	}
	return resUser, nil

}

func (uu *userUsecase) Login(user model.User) (string, error) {
	storedUser := model.User{}
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil

}
