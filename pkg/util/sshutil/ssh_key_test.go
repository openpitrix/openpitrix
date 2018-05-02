package sshutil

import "testing"

func TestMakeSSHKeyPair(t *testing.T) {
	private, public, err := MakeSSHKeyPair("ssh-rsa")
	if err != nil {
		t.Errorf("Error: %+v", err)
	}
	if private == "" || public == "" {
		t.Errorf("Generate rsa ssh key failed")
	}

	private, public, err = MakeSSHKeyPair("ssh-dsa")
	if err != nil {
		t.Errorf("Error: %+v", err)
	}
	if private == "" || public == "" {
		t.Errorf("Generate dsa ssh key failed")
	}
}
