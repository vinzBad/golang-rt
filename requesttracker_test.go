package rt

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	apiURL   = "http://localhost:8080/REST/1.0/"
	user     = "root"
	password = "password"
)

func TestNew(t *testing.T) {

	var (
		invalidAPIURL            = "invalid"
		validAPIURLWithoutSuffix = "http://example.com/ERROR"
	)
	_, err := New(invalidAPIURL, "", "")

	if err != ErrInvalidAPIURL {
		t.Errorf("New(%q) didn't detect invalid URL", invalidAPIURL)
	}

	_, err = New(validAPIURLWithoutSuffix, "", "")

	if err != ErrInvalidAPIURL {
		t.Errorf("New(%q) didn't detect missing api suffix in URL", validAPIURLWithoutSuffix)
	}

	_, err = New(apiURL, "", "")

	if err != nil {
		t.Errorf("New(%q), didn't accept valid api URL", apiURL)
	}
}

func TestLogin(t *testing.T) {
	// test invalid credentials handling
	tracker, err := New(apiURL, user, "")
	if err != nil {
		t.Errorf("Failed to initialize RT client: %q", err)
	}
	err = tracker.Login()

	if err == nil {
		t.Errorf("RT client logged in with invalid credentials")
	}

	tracker, err = New(apiURL, user, password)
	if err != nil {
		t.Errorf("Failed to initialize RT client: %q", err)
	}

	err = tracker.Login()
	if err != nil {
		t.Errorf("Login() with valid credentials failed: %q", err)
	}

	if tracker.isLoggedIn == false {
		t.Errorf("tracker.isLoggedIn is false, want true")
	}
}

func TestGetTicket(t *testing.T) {
	tracker, err := New(apiURL, user, password)
	if err != nil {
		t.Errorf("Failed to initialize RT client: %q", err)
	}
	err = tracker.Login()
	if err != nil {
		t.Errorf("Login() with valid credentials failed: %q", err)
	}
	expectedTicket := &Ticket{
		ID:         13,
		Queue:      "General",
		Owner:      "",
		Creator:    "root",
		Subject:    "test",
		Status:     "new",
		Priority:   "0",
		Requestors: []string{"root@localhost", "support@localhost"},
		Cc:         []string{""},
		AdminCc:    []string{""},
	}
	ticket, err := tracker.GetTicket(13)
	if err != nil {
		t.Errorf("GetTicket(13) failed: %q", err)
	}
	if !cmp.Equal(ticket, expectedTicket) {
		t.Errorf("GetTicket(13) didn't return expected data")
	}
}
