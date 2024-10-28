package handler

import (
	"encoding/json"

	"github.com/cvicream/cv-icream-api/auth"
	"github.com/cvicream/cv-icream-api/model"
	"github.com/cvicream/cv-icream-api/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func GoogleAuth(c *fiber.Ctx) error {
	redirectUrl := c.Query("redirect")
	log.Infof("Redirect URL: %s\n", redirectUrl)
	path := auth.ConfigGoogle()
	url := path.AuthCodeURL(redirectUrl)
	return c.Redirect(url)
}

func GoogleCallback(c *fiber.Ctx) error {
	// get redirect url from state
	redirectUrl := c.FormValue("state")
	log.Infof("Redirect URL: %s\n", redirectUrl)

	token, err := auth.ConfigGoogle().Exchange(c.Context(), c.FormValue("code"))
	if err != nil {
		log.Errorf("Could not get token: %s\n", err.Error())
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "could not get token", "data": nil})
	}

	client := auth.ConfigGoogle().Client(c.Context(), token)
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Errorf("Could not get user info: %s\n", err.Error())
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "could not get user info", "data": nil})
	}

	defer userInfo.Body.Close()

	var googleUserInfo GoogleUserInfo
	err = json.NewDecoder(userInfo.Body).Decode(&googleUserInfo)
	if err != nil {
		log.Errorf("Could not decode user info: %s\n", err.Error())
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "could not decode user info", "data": nil})
	}

	user, err := service.GetUserByEmail(googleUserInfo.Email, "google")
	if user == nil || err != nil {
		user, err = service.CreateUser(&model.User{
			FirstName:  &googleUserInfo.GivenName,
			LastName:   &googleUserInfo.FamilyName,
			Email:      googleUserInfo.Email,
			Avatar:     service.ConvertImageToBase64(googleUserInfo.Picture),
			Provider:   "google",
			ProviderID: &googleUserInfo.Sub,
		})

		if user == nil {
			log.Errorf("Could not create user: %s\n", err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	jwtToken, err := service.CreateUserJwtToken(user)
	if err != nil {
		log.Errorf("Could not create jwt token: %s\n", err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if err := service.CheckIfSurveyExist(user.ID); err != nil {
		if err.Error() == "user already has a survey" {
			return c.Redirect(redirectUrl + "/dashboard?token=" + jwtToken)
		}
	}

	return c.Redirect(redirectUrl + "/survey?token=" + jwtToken)
}
