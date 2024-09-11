package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cvicream/cv-icream-api/auth"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

func LinkedInAuth(c *fiber.Ctx) error {
	url := auth.ConfigLinkedIn().AuthCodeURL("state", oauth2.AccessTypeOffline)
	return c.Redirect(url)
}

func LinkedInCallback(c *fiber.Ctx) error {
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
	resp, err := client.Get("https://api.linkedin.com/v2/me")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get user info")
	}
	defer resp.Body.Close()

	var profile map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to decode user profile")
	}

	// Display the user profile
	return c.JSON(profile)
}
