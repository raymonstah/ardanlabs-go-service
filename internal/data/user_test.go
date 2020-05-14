package data_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/raymonstah/ardanlabs-go-service/internal/data"
	"github.com/raymonstah/ardanlabs-go-service/internal/platform/auth"
	"github.com/raymonstah/ardanlabs-go-service/internal/tests"
)

func TestUserCreate(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	t.Cleanup(teardown)

	ctx := context.Background()
	now := time.Now()
	user := data.NewUser{
		Email: "rho@launchdarkly.com",
		Name:  "Raymond Ho",
	}

	t.Log("creating user...")
	userCreated, err := data.Users.Create(ctx, db, user, now)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("retreiving user...")
	// claims is information about the person making the request.
	claims := auth.NewClaims(
		"718ffbea-f4a1-4667-8ae3-b349da52675e", // This is just some random UUID.
		[]string{auth.RoleAdmin, auth.RoleUser},
		now, time.Hour,
	)
	gotUser, err := data.Users.Retrieve(ctx, claims, db, userCreated.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(userCreated, gotUser); diff != "" {
		t.Fatal("retreived user not same as created user")
	}

}
