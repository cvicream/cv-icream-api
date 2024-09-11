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

func GetCVs() ([]model.CV, error) {
	cvs := []model.CV{}
	result := database.DB.Find(&cvs)
	if result.Error != nil {
		log.Errorf("Failed to retrieve CVs: %v", result.Error)
		return nil, result.Error
	}
	return cvs, nil
}

func GetCV(id string) (*model.CV, error) {
	cv := model.CV{}
	result := database.DB.Where("id = ?", id).First(&cv)
	if result.Error != nil {
		log.Errorf("Failed to retrieve CV: %v", result.Error)
		return nil, result.Error
	}
	return &cv, nil
}

func UpdateCV(cv *model.CV) (*model.CV, error) {
	result := database.DB.Save(&cv)
	if result.Error != nil {
		log.Errorf("Failed to update CV: %v", result.Error)
		return nil, result.Error
	}
	return cv, nil
}

func DeleteCV(id string) error {
	result := database.DB.Where("id = ?", id).Delete(&model.CV{})
	if result.Error != nil {
		log.Errorf("Failed to delete CV: %v", result.Error)
		return result.Error
	}
	return nil
}
