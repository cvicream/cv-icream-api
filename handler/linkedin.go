package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cvicream/cv-icream-api/auth"
	"github.com/cvicream/cv-icream-api/model"
	"github.com/cvicream/cv-icream-api/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type LinkedInUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func LinkedInAuth(c *fiber.Ctx) error {
	redirectUrl := c.Query("redirect")
	log.Infof("Redirect URL: %s\n", redirectUrl)
	path := auth.ConfigLinkedIn()
	url := path.AuthCodeURL(redirectUrl)
	return c.Redirect(url)
}

func LinkedInCallback(c *fiber.Ctx) error {
	// get redirect url from state
	redirectUrl := c.FormValue("state")
	log.Infof("Redirect URL: %s\n", redirectUrl)

	// Get the authorization code from LinkedIn
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing 'code' in query")
	}

	// Exchange code for OAuth2 token
	token, err := auth.ConfigLinkedIn().Exchange(context.Background(), code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to exchange token")
	}

	// Extract OpenID Connect ID token (JWT)
	idToken := token.Extra("id_token")
	if idToken == nil {
		return c.Status(fiber.StatusInternalServerError).SendString("No ID token found")
	}

	// Decode the ID Token (JWT)
	fmt.Printf("ID Token: %v\n", idToken)

	// Make request to LinkedIn API to get user profile
	client := auth.ConfigLinkedIn().Client(context.Background(), token)
	userInfo, err := client.Get("https://api.linkedin.com/v2/userinfo")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get user info")
	}
	defer userInfo.Body.Close()

	var linkedinUserInfo LinkedInUserInfo
	err = json.NewDecoder(userInfo.Body).Decode(&linkedinUserInfo)
	if err != nil {
		log.Errorf("Could not decode user info: %s\n", err.Error())
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "could not decode user info", "data": nil})
	}

	user, err := service.GetUserByEmail(linkedinUserInfo.Email, "linkedin")
	if user == nil || err != nil {
		user, err = service.CreateUser(&model.User{
			FirstName:  &linkedinUserInfo.GivenName,
			LastName:   &linkedinUserInfo.FamilyName,
			Email:      linkedinUserInfo.Email,
			Avatar:     service.ConvertImageToBase64(linkedinUserInfo.Picture),
			Provider:   "linkedin",
			ProviderID: &linkedinUserInfo.Sub,
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
