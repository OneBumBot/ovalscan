package ssh

import (
	"errors"

	"golang.org/x/crypto/ssh"
)

func CreateSession(client *ssh.Client) (*ssh.Session, error) {

	session, err := client.NewSession()
	if err != nil {
		return nil, errors.New("failed to create ssh session")
	}

	return session, nil
}

