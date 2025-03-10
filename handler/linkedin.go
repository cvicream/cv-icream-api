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
	signInUrl := redirectUrl + "/sign-in"

	// Get the authorization code from LinkedIn
	code := c.Query("code")
	if code == "" {
		log.Infof("Missing 'code' in query")
		return c.Redirect(signInUrl)
	}

	// Exchange code for OAuth2 token
	token, err := auth.ConfigLinkedIn().Exchange(context.Background(), code)
	if err != nil {
		log.Infof("Failed to exchange token")
		return c.Redirect(signInUrl)
	}

	// Extract OpenID Connect ID token (JWT)
	idToken := token.Extra("id_token")
	if idToken == nil {
		log.Infof("No ID token found")
		return c.Redirect(signInUrl)
	}

	// Decode the ID Token (JWT)
	fmt.Printf("ID Token: %v\n", idToken)

	// Make request to LinkedIn API to get user profile
	client := auth.ConfigLinkedIn().Client(context.Background(), token)
	userInfo, err := client.Get("https://api.linkedin.com/v2/userinfo")
	if err != nil {
		log.Infof("Failed to get user info")
		return c.Redirect(signInUrl)
	}
	defer userInfo.Body.Close()

	var linkedinUserInfo LinkedInUserInfo
	err = json.NewDecoder(userInfo.Body).Decode(&linkedinUserInfo)
	if err != nil {
		log.Errorf("Could not decode user info: %s\n", err.Error())
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "could not decode user info", "data": nil})
	}

	user, err := service.GetUserByEmail(linkedinUserInfo.Email, "linkedin")
	isNewUser := false
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

		isNewUser = true
	}

	jwtToken, err := service.CreateUserJwtToken(user)
	if err != nil {
		log.Errorf("Could not create jwt token: %s\n", err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	errSurvey := service.CheckIfSurveyExist(user.ID)
	if (errSurvey != nil && errSurvey.Error() == "user already has a survey") || !isNewUser {
		return c.Redirect(redirectUrl + "/dashboard?token=" + jwtToken)
	}

	return c.Redirect(redirectUrl + "/survey?token=" + jwtToken)
}
