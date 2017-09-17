# flick/session

[![](https://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/olafal0/flick/session)

This subpackage is designed to provide a basic and straightforward way to create
and manage HTTP sessions using cookies. Sessions are stored in-memory and
safe for concurrent access. Methods are also provided to save the whole session
list to disk in [gob](https://golang.org/pkg/encoding/gob/) format. Since
incremental saving is not at all supported, this method is unsuitable for
anything beyond a small number of active sessions.

# Usage

Create a new session for a given email (valid for 30 days by default):

`sessionCookie, emailCookie := session.BeginSession(userEmail)`

Set the cookies in the response writer when the user first logs in:

`http.SetCookie(c.Wr, sessionCookie)`

`http.SetCookie(c.Wr, emailCookie)`

Now, you can easily validate that a user is who they say they are on any request:

`isValid, validatedEmail := session.VerifyRequest(c.Req)`

If persistence is needed, there's an easy (albeit limited) way to get it.
This starts a goroutine that saves every 60 seconds (saving and loading are
thread-safe):

`session.StartAutosaving("sessions.gob", time.Minute)`

At startup, you can load your sessions back:

`session.Load("sessions.gob")`

To prevent data loss (for example, from exiting mid-save), you can freeze the
sessions list: it will not be modified or saved until you exit, and this call
will block until all saving is finished. Don't call this unless you are for sure
going to exit (it takes a write lock on the sessions list and never releases it).

`session.Freeze()`