package gcp

import (
	"context"
	"fmt"
	"time"

	"github.com/nimil-jp/gin-utils/errors"

	"go-gin-ddd/config"
	"go-gin-ddd/driver/gcp"
)

type IFirebase interface {
	AuthClient() gcp.FirebaseAuthClient
	SetClaimsUID(firebaseUID string, uid uint) error
}

type firebase struct {
	client gcp.FirebaseAuthClient
}

func NewFirebase() IFirebase {
	return firebase{
		client: gcp.AuthClient(),
	}
}

func (i firebase) AuthClient() gcp.FirebaseAuthClient {
	return i.client
}

func (i firebase) SetClaimsUID(firebaseUID string, uid uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	user, err := i.client.GetUser(ctx, firebaseUID)
	if err != nil {
		return errors.NewUnexpected(err)
	}

	if user.CustomClaims == nil {
		user.CustomClaims = map[string]interface{}{}
	}

	user.CustomClaims[fmt.Sprintf("%s_id", config.AppName)] = uid

	err = i.client.SetCustomUserClaims(ctx, firebaseUID, user.CustomClaims)
	if err != nil {
		return errors.NewUnexpected(err)
	}

	return nil
}
