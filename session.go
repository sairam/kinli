package kinli

import (
	"net/http"
	"reflect"

	"github.com/gorilla/sessions"
)

var (
	// SessionName ..
	SessionName = "_kinli" // can be seen in the cookies list

	// HomePathNonAuthed redirection for NonLoggedInUser
	HomePathNonAuthed = "/home"

	// HomePathAuthed redirection for LoggedInUser
	HomePathAuthed = "/"

	// SessionStore interface to will be used
	SessionStore sessions.Store
)

/*

Examples

to override session.Options:
session.Options = &sessions.Options{
    Path:     "/",
    Domain: "domain.com",
    MaxAge:   86400 * 7, // 7 days
    HttpOnly: true,
    Secure: true,
}

cookie store:
SESSION_AUTHENTICATION_KEY="<string of any length>"
sessionStore = sessions.NewCookieStore(
        []byte(os.Getenv("SESSION_AUTHENTICATION_KEY")),
        ek, // optional encryption key of 32 bytes
    )
// refer https://github.com/heroku-examples/go-sessions-demo/blob/2b6ba688c181dc15f9898c03754c2f1d8e85cd48/main.go#L93-L107

file system store:
sessions.NewFilesystemStore("./sessions", []byte(os.Getenv("SESSION_AUTHENTICATION_KEY")))

redis store:
import redisStore "gopkg.in/boj/redistore.v1"
sessionStore = redisStore.NewRediStore(5, "tcp", ":6379", "redis-password",
  []byte(os.Getenv("SESSION_AUTHENTICATION_KEY")),
  ek,
)

*/

// HttpContext .. is the current context
type HttpContext struct {
	W http.ResponseWriter
	R *http.Request
}

// RedirectAfterAuth will redirect user immediately after authentication
// returns true if redirected, false if not
func (hc *HttpContext) RedirectAfterAuth(flash string) bool {
	if hc.isAuthed() {
		session := hc.getSession()
		if session == nil {
			return false
		}
		if flash != "" {
			session.AddFlash(flash)
			hc.saveSession(session)
		}

		http.Redirect(hc.W, hc.R, HomePathAuthed, 302)
		return true
	}
	return false
}

// RedirectUnlessAuthed redirects if the user is not logged in
// returns true if redirected, false if not
func (hc *HttpContext) RedirectUnlessAuthed(flash string) bool {
	if !hc.isAuthed() {
		session := hc.getSession()
		if session == nil {
			return true
		}
		if flash != "" {
			session.AddFlash(flash)
			hc.saveSession(session)
		}

		http.Redirect(hc.W, hc.R, HomePathNonAuthed, 302)
		return true
	}
	return false
}

// Helper methods for the session handlers
func (hc *HttpContext) getSession() *sessions.Session {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, err := SessionStore.Get(hc.R, SessionName)
	if err != nil {
		http.Error(hc.W, err.Error(), http.StatusInternalServerError)
		return nil
	}
	return session
}

func (hc *HttpContext) saveSession(session *sessions.Session) {
	// Save it before we write to the response/return from the handler.
	session.Save(hc.R, hc.W)
}

func (hc *HttpContext) clearFlashes() {
	session := hc.getSession()
	if flashes := session.Flashes(); len(flashes) > 0 {
		hc.saveSession(session)
	}
}

// GetFlashes gets the list of flashes and flushes as well
func (hc *HttpContext) GetFlashes() []string {
	session := hc.getSession()
	if flashes := session.Flashes(); len(flashes) > 0 {
		fs := make([]string, len(flashes))
		for i, f := range flashes {
			fs[i] = string(reflect.ValueOf(f).String())
		}
		hc.saveSession(session)
		return fs
	}
	return nil
}

func (hc *HttpContext) AddFlash(flash string) {

	session := hc.getSession()
	if session == nil {
		return
	}

	session.AddFlash(flash)
	hc.saveSession(session)
}

// setSessionData sets a key/value into a session and saves it
// value should be registered first with gob.Register(&Data{})
// should be sending data should be of the format ("data", &Data{})
// retrieving should be of the format data, ok := val.(*Data)
func (hc *HttpContext) SetSessionData(key string, value interface{}) {
	session := hc.getSession()
	if session == nil {
		return
	}

	session.Values[key] = value
	hc.saveSession(session)
}

// TODO test this if this should be a pointer
// func (hc *HttpContext) GetSessionIntoData(key string, value interface{}) (ok bool) {
// 	ifaceType := reflect.TypeOf(value) // see gob.encoding how to reflect pointers
// 	v, ok = val.(ifaceType)
// 	*value = v
// 	return
// }

// GetSessionData can be used to retrieve the data
// retrieving should be of the format data, ok := val.(*Data)
func (hc *HttpContext) GetSessionData(key string) interface{} {
	session := hc.getSession()
	if session == nil {
		return nil
	}
	return session.Values[key]
}

// ClearSession logs out the user by deleting all session data
func (hc *HttpContext) ClearSession() {
	session := hc.getSession()
	if session == nil {
		return
	}

	for key := range session.Values {
		delete(session.Values, key)
	}

	hc.saveSession(session)
}

func (hc *HttpContext) isAuthed() bool {
	return IsAuthed(hc)
}

// IsAuthed is used to verify if a user is logged in or not
// can be overriden with your own business logic
var IsAuthed = isAuthed

func isAuthed(hc *HttpContext) bool {
	if len(hc.getSession().Values) > 0 {
		return true
	}
	return false
}
