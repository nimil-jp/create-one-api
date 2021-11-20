package gcp

import (
	"context"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"

	"go-gin-ddd/config"
)

var (
	gcsClient *storage.Client
)

func init() {
	var err error

	ctx := context.Background()
	gcsClient, err = storage.NewClient(ctx, option.WithCredentialsFile(config.Env.GCP.CredentialPath))
	if err != nil {
		panic(err)
	}
}

func GcsClient() *storage.Client {
	return gcsClient
}
