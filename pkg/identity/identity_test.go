package identity

import "testing"

func TestGetIdentity(t *testing.T) {
	identity := getIdentity()
	t.Log(identity)
}
