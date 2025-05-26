package auth

import (
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"log"
	"os"
)

type Environment struct {
	host         string
	clientID     string
	clientSecret string
	callback     string
}

type Auth struct {
	Provider        *oidc.Provider
	Config          *oauth2.Config
	IDTokenVerifier *oidc.IDTokenVerifier
}

// Read environment vars for configuring the
// Oauth2 client
func ReadEnvironment() *Environment {

	env := &Environment{
		host:         os.Getenv("AUTH_SERVER"),
		clientID:     os.Getenv("CLIENT_ID"),
		clientSecret: os.Getenv("CLIENT_SECRET"),
		callback:     os.Getenv("CALLBACK"),
	}

	if env.host == "" {
		fmt.Println("AUTH_SERVER has not been set")
	}
	if env.clientID == "" {
		fmt.Println("CLIENT_ID has not been set")
	}
	if env.clientSecret == "" {
		fmt.Println("CLIENT_SECRET has not been set")
	}
	if env.callback == "" {
		fmt.Println("CALLBACK has not been set")
	}

	log.Println(env.host)
	log.Println(env.clientID)
	log.Println(env.clientSecret)
	log.Println(env.callback)

	return env
}

func NewAuth() *Auth {

	env := ReadEnvironment()

	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, env.host)
	if err != nil {
		log.Fatal(err)
	}

	clientID := env.clientID
	clientSecret := env.clientSecret

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  env.callback,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}

	verifier := provider.Verifier(oidcConfig)

	return &Auth{
		Provider:        provider,
		Config:          config,
		IDTokenVerifier: verifier,
	}
}
