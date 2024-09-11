package auth

import (
	"github.com/cvicream/cv-icream-api/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
)

func ConfigLinkedIn() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     config.Config("LINKEDIN_CLIENT_ID"),
		ClientSecret: config.Config("LINKEDIN_CLIENT_SECRET"),
		RedirectURL:  config.Config("LINKEDIN_REDIRECT_URL"),
		Scopes:       []string{"r_liteprofile", "r_emailaddress", "openid"},
		Endpoint:     linkedin.Endpoint,
	}

	return conf
}
