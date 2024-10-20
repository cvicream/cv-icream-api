package service

import (
	"encoding/base64"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2/log"
)

func ConvertImageToBase64(imageUrl string) string {
	// Download the image
	response, err := http.Get(imageUrl)
	if err != nil {
		log.Errorf("Error downloading image: %v", err)
		return ""
	}
	defer response.Body.Close()

	// Read the image bytes
	imageBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorf("Error reading image bytes: %v", err)
		return ""
	}

	// Encode to base64
	base64String := base64.StdEncoding.EncodeToString(imageBytes)
	return base64String
}
