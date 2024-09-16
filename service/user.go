package service

import (
	"time"

	"github.com/cvicream/cv-icream-api/config"
	"github.com/cvicream/cv-icream-api/database"
	"github.com/cvicream/cv-icream-api/model"

	"github.com/golang-jwt/jwt"
)

func CreateUserJwtToken(user *model.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString([]byte(config.Config("JWT_SECRET")))
	return tokenString, err
}

func GetUserById(id float64) (*model.User, error) {
	var user model.User
	result := database.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func CreateUser(user *model.User) (*model.User, error) {
	result := database.DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func UpdateUser(id float64, user *model.User) (*model.User, error) {
	var existingUser model.User
	if err := database.DB.Where("id = ?", id).First(&existingUser).Error; err != nil {
		return nil, err
	}

	existingUser.FirstName = user.FirstName
	existingUser.LastName = user.LastName

	result := database.DB.Save(&existingUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return &existingUser, nil
}
