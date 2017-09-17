package user

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	SetHashCost(10)

	u := AuthUser{}
	err := u.SetPassword("clowndentists")
	if err != nil {
		t.Error("Could not set password")
	}

	if !u.Authenticate("clowndentists") {
		t.Error("Could not authenticate with correct password")
	}

	if u.Authenticate("dentistclowns") {
		t.Error("Could authenticate with incorrect password")
	}
}

func TestNewExtendedUser(t *testing.T) {
	SetHashCost(10)

	type ExtendedUser struct {
		AuthUser
		Name  string
		Email string
	}
	eu := ExtendedUser{Name: "Todd Chavez"}
	eu.Email = "its.todd.its.me@yahoo.com"
	err := eu.SetPassword("clowndentists")
	if err != nil {
		t.Error("Could not set password for extended user")
	}

	if !eu.Authenticate("clowndentists") {
		t.Error("Could not authenticate with correct password for extended user")
	}

	if eu.Authenticate("dentistclowns") {
		t.Error("Could authenticate with incorrect password for extended user")
	}
}

func BenchmarkSetPassword(b *testing.B) {
	u := AuthUser{}
	for n := 0; n < b.N; n++ {
		u.SetPassword("password" + string(n))
	}
}
