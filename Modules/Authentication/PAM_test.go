package Authentication

import (
	"testing"
)

func TestExist(t *testing.T) {
	r := UserExist("fransebas")
	if !r {
		t.Fail()
	}

	r = UserExist("frans3baz")

	if r {
		t.Fail()
	}
}

func TestAuth(t *testing.T) {

	password := "" // right password here

	// valid
	r, _ := ValidateUser("fransebas", password)

	if !r {
		t.Fail()
	}

	// wrong password
	r, _ = ValidateUser("fransebas", "abc123")
	if r {
		t.Fail()
	}

	// wrong password and username
	r, _ = ValidateUser("fr4ns3baz", "abc123")
	if r {
		t.Fail()
	}

	// wrong username
	r, _ = ValidateUser("fr4ns3baz", password)
	if r {
		t.Fail()
	}
}
