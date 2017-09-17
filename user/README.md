# flick/user

[![](https://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/olafal0/flick/user)

The user subpackage provides the AuthUser type, which is designed to add
authentication functionality to your own User types using composition.

# Usage

In any custom user type, embed the AuthUser type:

```go
type MyUser struct {
    user.AuthUser
    Name string
    Email string
}
```

This gets you basic authentication functionality:

```go
u := MyUser{ Name: "Todd Chavez", Email: "its.todd.its.me@yahoo.com" }
u.SetPassword("abc123")
u.Authenticate("abc123") // returns true
```

Additionally, marshaling `MyUser` will include the password (hashed and salted
using bcrypt) so it is suitable for saving without any modifications.