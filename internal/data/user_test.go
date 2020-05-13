package data_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/raymonstah/ardanlabs-go-service/internal/data"
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
	gotUser, err := data.Users.Retrieve(ctx, db, userCreated.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(userCreated, gotUser); diff != "" {
		t.Fatal("retreived user not same as created user")
	}

}
