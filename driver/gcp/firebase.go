package gcp

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

	"go-gin-ddd/config"
)

var authClient *auth.Client

func init() {
	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsFile(config.Env.GCP.CredentialPath))
	if err != nil {
		panic(err)
	}

	authClient, err = app.Auth(context.Background())
	if err != nil {
		panic(err)
	}
}

type FirebaseAuthClient interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
	GetUser(ctx context.Context, uid string) (*auth.UserRecord, error)
	SetCustomUserClaims(ctx context.Context, uid string, customClaims map[string]interface{}) error
}

func AuthClient() FirebaseAuthClient {
	return authClient
}
