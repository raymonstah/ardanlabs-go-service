package tests

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/raymonstah/ardanlabs-go-service/internal/data"
	"github.com/raymonstah/ardanlabs-go-service/internal/platform/auth"
	"github.com/raymonstah/ardanlabs-go-service/internal/platform/database"
	"github.com/raymonstah/ardanlabs-go-service/internal/platform/database/databasetest"
	"github.com/raymonstah/ardanlabs-go-service/internal/platform/web"
	"github.com/raymonstah/ardanlabs-go-service/internal/schema"
)

// Success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

// NewUnit creates a test database inside a Docker container. It creates the
// required table structure but the database is otherwise empty.
//
// It does not return errors as this intended for testing only. Instead it will
// call Fatal on the provided testing.T if anything goes wrong.
//
// It returns the database to use as well as a function to call at the end of
// the test.
func NewUnit(t *testing.T) (*sqlx.DB, func()) {
	t.Helper()

	c := databasetest.StartContainer(t)

	db, err := database.Open(database.Config{
		User:       "postgres",
		Password:   "postgres",
		Host:       c.Host,
		Name:       "postgres",
		DisableTLS: true,
	})
	if err != nil {
		t.Fatalf("opening database connection: %v", err)
	}

	t.Log("waiting for database to be ready")

	// Wait for the database to be ready. Wait 100ms longer between each attempt.
	// Do not try more than 20 times.
	var pingError error
	maxAttempts := 20
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		pingError = db.Ping()
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	}

	if pingError != nil {
		databasetest.DumpContainerLogs(t, c)
		databasetest.StopContainer(t, c)
		t.Fatalf("waiting for database to be ready: %v", pingError)
	}

	if err := schema.Migrate(db); err != nil {
		databasetest.StopContainer(t, c)
		t.Fatalf("migrating: %s", err)
	}

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		db.Close()
		databasetest.StopContainer(t, c)
	}

	return db, teardown
}

// Test owns state for running and shutting down tests.
type Test struct {
	DB            *sqlx.DB
	Log           *log.Logger
	Authenticator *auth.Authenticator

	t       *testing.T
	cleanup func()
}

// NewIntegration creates a database, seeds it, constructs an authenticator.
func NewIntegration(t *testing.T) *Test {
	t.Helper()

	// Initialize and seed database. Store the cleanup function call later.
	db, cleanup := NewUnit(t)

	if err := data.Seed(db); err != nil {
		t.Fatal(err)
	}

	// Create the logger to use.
	logger := log.New(os.Stdout, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// Create RSA keys to enable authentication in our service.
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	// Build an authenticator using this static key.
	kid := "4754d86b-7a6d-4df5-9c65-224741361492"
	kf := auth.NewSimpleKeyLookupFunc(kid, key.Public().(*rsa.PublicKey))
	authenticator, err := auth.NewAuthenticator(key, kid, "RS256", kf)
	if err != nil {
		t.Fatal(err)
	}

	return &Test{
		DB:            db,
		Log:           logger,
		Authenticator: authenticator,
		t:             t,
		cleanup:       cleanup,
	}
}

// Teardown releases any resources used for the test.
func (test *Test) Teardown() {
	test.cleanup()
}

// Token generates an authenticated token for a user.
func (test *Test) Token(email, pass string) string {
	test.t.Helper()

	claims, err := data.Users.Authenticate(
		context.Background(), test.DB, time.Now(),
		email, pass,
	)
	if err != nil {
		test.t.Fatal(err)
	}

	tkn, err := test.Authenticator.GenerateToken(claims)
	if err != nil {
		test.t.Fatal(err)
	}

	return tkn
}

// Context returns an app level context for testing.
func Context() context.Context {
	values := web.Values{
		TraceID: uuid.New().String(),
		Now:     time.Now(),
	}

	return context.WithValue(context.Background(), web.KeyValues, &values)
}
