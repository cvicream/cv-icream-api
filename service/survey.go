package service

import (
	"errors"
	"fmt"

	"github.com/cvicream/cv-icream-api/database"
	"github.com/cvicream/cv-icream-api/model"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

func CheckIfSurveyExist(id uint) error {
	tx := database.DB.Begin()
	// Check if user already has a survey
	var existingSurvey model.Survey
	err := tx.Where("user_id = ?", id).First(&existingSurvey).Error
	if err == nil {
		tx.Rollback()
		return errors.New("user already has a survey")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return errors.New(fmt.Sprintf("failed to check existing survey: %v", err))
	}
	return nil
}

func CreateSurvey(survey *model.Survey) error {
	// Start transaction
	tx := database.DB.Begin()

	// Check if user exists
	var user model.User
	if err := tx.First(&user, survey.UserID).Error; err != nil {
		tx.Rollback()
		return errors.New("user not found")
	}

	if err := CheckIfSurveyExist(survey.UserID); err != nil {
		return err
	}

	// Create survey
	if err := tx.Create(survey).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	return tx.Commit().Error
}

func DeleteSurveys(userId float64) error {
	result := database.DB.Unscoped().Where("user_id = ?", userId).Delete(&model.Survey{})
	if result.Error != nil {
		log.Errorf("Failed to delete surveys: %v", result.Error)
		return result.Error
	}
	return nil
}
