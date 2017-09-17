package session

import (
	"net/http"
	"os"
	"testing"
)

func TestSessionCreation(t *testing.T) {
	// create a user and session
	userEmail := "123@abc.com"
	sessionCookie, emailCookie := BeginSession(userEmail)

	if emailCookie.Value != userEmail {
		t.Error("BeginSession() set incorrect email cookie")
	}

	if !Verify(emailCookie.Value, sessionCookie.Value) {
		t.Error("Valid session could not be verified")
	}

	if Verify(emailCookie.Value, "") {
		t.Error("Empty session key wrongly resulted in verified session")
	}

	if Verify(emailCookie.Value+"a", sessionCookie.Value) {
		t.Error("Different email wrongly resulted in verified session")
	}

	r, _ := http.NewRequest("get", "https://localhost/test", nil)
	r.AddCookie(sessionCookie)
	r.AddCookie(emailCookie)

	if valid, _ := VerifyRequest(r); !valid {
		t.Error("Valid request gave incorrect result")
	}

	// force expire session cookie
	session := getSession(userEmail)
	session.Created = session.Created.Add(-2 * SessionDuration)

	if valid, _ := VerifyRequest(r); valid {
		t.Error("Expired session did not invalidate session")
	}

	// un-expire session and check that it was deleted
	session.Created = session.Created.Add(2 * SessionDuration)

	if valid, _ := VerifyRequest(r); valid {
		t.Error("Expired session did not get deleted")
	}

	// make a new session again, and validate that we can delete it
	sessionCookie, emailCookie = BeginSession(userEmail)
	r, _ = http.NewRequest("get", "https://localhost/test", nil)
	r.AddCookie(sessionCookie)
	r.AddCookie(emailCookie)
	if valid, _ := VerifyRequest(r); !valid {
		t.Error("Could not validate re-added session")
	}

	Save("test_sessions.gob")
	// overwrite everything
	currentSessions = make(map[string]*Session)
	// verify that our valid session isn't there anymore
	if valid, _ := VerifyRequest(r); valid {
		t.Error("Session should have been removed")
	}
	// load our sessions back (and clean up the file)
	Load("test_sessions.gob")
	os.Remove("test_sessions.gob")

	if valid, _ := VerifyRequest(r); !valid {
		t.Error("Session should have been loaded back")
	}

	EndSession(userEmail)
	if valid, _ := VerifyRequest(r); valid {
		t.Error("Ended session did not invalidate it correctly")
	}
}
