package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOauthConfig *oauth2.Config

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, proceeding without it")
	} else {
		fmt.Println(".env file loaded successfully")
	}
	GoogleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// GetGoogleOAuthURL generates the Google OAuth2 URL with state parameter
func GetGoogleOAuthURL(state string) string {
	return GoogleOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

// ExchangeCode exchanges the authorization code for an OAuth2 token
func ExchangeCode(code string) (*oauth2.Token, error) {
	token, err := GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	return token, nil
}
