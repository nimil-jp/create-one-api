package gcp

import (
	"context"
	"fmt"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/nimil-jp/gin-utils/errors"

	"go-gin-ddd/config"
	"go-gin-ddd/driver/gcp"
)

type IFirebase interface {
	AuthClient() *auth.Client
	SetClaimsUID(firebaseUID string, uid uint) error
}

type firebase struct {
	client *auth.Client
}

func NewFirebase() IFirebase {
	return firebase{
		client: gcp.AuthClient(),
	}
}

func (i firebase) AuthClient() *auth.Client {
	return i.client
}

func (i firebase) SetClaimsUID(firebaseUID string, uid uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	user, err := i.client.GetUser(ctx, firebaseUID)
	if err != nil {
		return errors.NewUnexpected(err)
	}

	user.CustomClaims[fmt.Sprintf("%s_id", config.AppName)] = uid

	return i.client.SetCustomUserClaims(ctx, firebaseUID, user.CustomClaims)
}
