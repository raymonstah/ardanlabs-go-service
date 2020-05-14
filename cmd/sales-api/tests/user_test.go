package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/raymonstah/ardanlabs-go-service/cmd/sales-api/internal/handlers"
	"github.com/raymonstah/ardanlabs-go-service/internal/tests"
)

// TestUsers is the entry point for testing user management functions.
func TestUsers(t *testing.T) {
	test := tests.NewIntegration(t)
	t.Cleanup(test.Teardown)

	shutdown := make(chan os.Signal, 1)
	tests := UserTests{
		app:        handlers.API("develop", shutdown, test.Log, test.DB, test.Authenticator),
		userToken:  test.Token("user@example.com", "gophers"),
		adminToken: test.Token("admin@example.com", "gophers"),
	}

	t.Run("getUser400", tests.getUser400)

}

// UserTests holds methods for each user subtest. This type allows passing
// dependencies for tests while still providing a convenient syntax when
// subtests are registered.
type UserTests struct {
	app        http.Handler
	userToken  string
	adminToken string
}

// getUser400 validates a user request for a malformed userid.
func (ut *UserTests) getUser400(t *testing.T) {
	id := "12345"

	r := httptest.NewRequest("GET", "/v1/users/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+ut.adminToken)

	ut.app.ServeHTTP(w, r)

	t.Log("Given the need to validate getting a user with a malformed userid.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new user %s.", testID, id)
		{
			if w.Code != http.StatusBadRequest {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 400 for the response : %v", tests.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 400 for the response.", tests.Success, testID)

			recv := w.Body.String()
			resp := `{"error":"ID is not in its proper form"}`
			if resp != recv {
				t.Logf("\t\tTest %d:\tGot : %v", testID, recv)
				t.Logf("\t\tTest %d:\tWant: %v", testID, resp)
				t.Fatalf("\t%s\tTest %d:\tShould get the expected result.", tests.Failed, testID)
			}
			t.Logf("\t%s\tTest %d:\tShould get the expected result.", tests.Success, testID)
		}
	}
}
