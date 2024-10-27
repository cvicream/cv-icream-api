package service

import (
	"github.com/cvicream/cv-icream-api/database"
	"github.com/cvicream/cv-icream-api/model"
	"github.com/gofiber/fiber/v2/log"
)

func CreateCV(cv model.CV) (*model.CV, error) {
	result := database.DB.Create(&cv)
	if result.Error != nil {
		log.Errorf("Failed to create CV: %v", result.Error)
		return nil, result.Error
	}
	return &cv, nil
}

func GetCVs(userId float64) ([]model.CV, error) {
	cvs := []model.CV{}
	result := database.DB.Where("user_id = ?", userId).Order("updated_at desc").Find(&cvs)
	if result.Error != nil {
		log.Errorf("Failed to retrieve CVs: %v", result.Error)
		return nil, result.Error
	}
	return cvs, nil
}

func GetCV(userId float64, id string) (*model.CV, error) {
	cv := model.CV{}
	result := database.DB.Where("user_id = ? AND id = ?", userId, id).First(&cv)
	if result.Error != nil {
		log.Errorf("Failed to retrieve CV: %v", result.Error)
		return nil, result.Error
	}
	return &cv, nil
}

func UpdateCV(userId float64, cv *model.CV) (*model.CV, error) {
	result := database.DB.Save(&cv)
	if result.Error != nil {
		log.Errorf("Failed to update CV: %v", result.Error)
		return nil, result.Error
	}
	return cv, nil
}

func UpdateCVTitle(userId float64, cv *model.CV) (*model.CV, error) {
	result := database.DB.Omit("updated_at").Save(&cv)
	if result.Error != nil {
		log.Errorf("Failed to update CV: %v", result.Error)
		return nil, result.Error
	}
	return cv, nil
}

func DeleteCV(userId float64, id string) error {
	result := database.DB.Where("user_id = ? AND id = ?", userId, id).Delete(&model.CV{})
	if result.Error != nil {
		log.Errorf("Failed to delete CV: %v", result.Error)
		return result.Error
	}
	return nil
}
