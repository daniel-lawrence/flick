package session

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/gob"
	"log"
	"math/big"
	"net/http"
	"os"
	"sync"
	"time"
)

var cryptoRandMax *big.Int

// SessionDuration represents how long a session is valid for. 30-day sessions
// are valid by default
const SessionDuration = time.Hour * 24 * 30

const sessionCookieName = "flick-session"
const identifierCookieName = "flick-session-ident"

// Session represents a saved user's session.
type Session struct {
	Created time.Time
	Key     string
}

// maps email addresses (i.e. user IDs) to sessions
var currentSessions map[string]*Session
var sessionLock sync.RWMutex

func init() {
	cryptoRandMax = big.NewInt(0)
	cryptoRandMax.Exp(big.NewInt(2), big.NewInt(256), nil)

	currentSessions = make(map[string]*Session)
}

// getSession gets the Session for the specified email address if it exists
// (or nil otherwise) in a thread-safe manner.
func getSession(email string) *Session {
	sessionLock.RLock()
	s, ok := currentSessions[email]
	sessionLock.RUnlock()
	if !ok {
		return nil
	}
	return s
}

// setSession sets a value in the session table in a thread-safe manner.
func setSession(email string, s *Session) *Session {
	sessionLock.Lock()
	currentSessions[email] = s
	sessionLock.Unlock()
	return s
}

func deleteSession(email string) {
	sessionLock.Lock()
	delete(currentSessions, email)
	sessionLock.Unlock()
}

// BeginSession generates a new session token for this user, stores it, and
// returns the cookies to give to the user
func BeginSession(userEmail string) (session *http.Cookie, email *http.Cookie) {
	key, err := rand.Int(rand.Reader, cryptoRandMax)
	if err != nil {
		log.Println(err)
		return
	}

	keyStr := base64.StdEncoding.EncodeToString(key.Bytes())

	setSession(userEmail, &Session{time.Now(), keyStr})

	session = &http.Cookie{
		Name:     sessionCookieName,
		Value:    keyStr,
		Expires:  time.Now().Add(SessionDuration),
		Secure:   true,
		HttpOnly: true,
	}
	email = &http.Cookie{
		Name:     identifierCookieName,
		Value:    userEmail,
		Expires:  time.Now().Add(SessionDuration),
		Secure:   true,
		HttpOnly: true,
	}

	return
}

// EndSession removes this session, and returns blank session and email cookies
// to overwrite the user's current cookies.
func EndSession(userEmail string) (session *http.Cookie, email *http.Cookie) {
	// construct expired blank cookies to replace the user's current cookies
	session = &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Expires:  time.Now(),
		Secure:   true,
		HttpOnly: true,
	}
	email = &http.Cookie{
		Name:     identifierCookieName,
		Value:    "",
		Expires:  time.Now(),
		Secure:   true,
		HttpOnly: true,
	}
	deleteSession(userEmail)
	return
}

// VerifyRequest wraps the Verify function, and takes a *http.Request for
// convenience. The email and session key cookies will be automatically
// extracted and used. For convenience, it returns the email string of the
// valid and verified user, or an empty string on failure.
func VerifyRequest(r *http.Request) (valid bool, email string) {
	email = ""

	sessionCookie, sessionErr := r.Cookie(sessionCookieName)
	emailCookie, emailErr := r.Cookie(identifierCookieName)

	if sessionErr != nil || emailErr != nil {
		return false, email
	}

	email = emailCookie.Value
	if Verify(emailCookie.Value, sessionCookie.Value) {
		return true, email
	}
	return false, email
}

// Verify verifies that a given email and session key are unexpired, valid, and
// matching. Returns true if valid and false otherwise.
func Verify(emailKey, sessionKey string) bool {
	if emailKey == "" || sessionKey == "" {
		return false
	}

	session := getSession(emailKey)

	if session == nil {
		// session not found
		return false
	}

	expirationTime := session.Created.Add(SessionDuration)
	if time.Now().After(expirationTime) {
		// session has expired, remove it
		deleteSession(emailKey)
		return false
	}

	// use crypto/subtle to compare in constant time, just in case
	match := subtle.ConstantTimeCompare(
		[]byte(session.Key),
		[]byte(sessionKey),
	)

	if match == 1 {
		// session valid
		return true
	}

	return false
}

// Save writes the current session list to disk as a gob.
func Save(filename string) error {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		return err
	}
	w := gob.NewEncoder(f)

	sessionLock.RLock()
	err = w.Encode(currentSessions)
	sessionLock.RUnlock()

	return err
}

// Load loads a session list from a gob file.
func Load(filename string) error {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return err
	}
	d := gob.NewDecoder(f)

	sessionLock.Lock()
	err = d.Decode(&currentSessions)
	sessionLock.Unlock()

	return err
}

// StartAutosaving repeatedly saves the session list at a given interval. This
// function is non-blocking, as it creates and runs its own thread. When called,
// it does not immediately perform a save; it waits for the given interval first.
func StartAutosaving(filename string, interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			Save(filename)
		}
	}()
}

// Freeze takes the session lock for writing, thereby preventing any further
// reads or writes. Use this immediately before exiting to guarantee that no
// data will be lost if exiting mid-save. Note that changes since the last save
// may still be lost.
//
// This function will block until all reads or writes to the session table are
// finished.
func Freeze() {
	sessionLock.Lock()
}
