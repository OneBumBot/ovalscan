package ssh

import (
	"errors"
	"strings"

	"golang.org/x/crypto/ssh"
)

func CreateSession(client *ssh.Client) (*ssh.Session, error) {

	session, err := client.NewSession()
	if err != nil {
		return nil, errors.New("failed to create ssh session")
	}

	return session, nil
}

func ExecuteCommands(session *ssh.Session, cmds ...string) ([]byte, error) {
	if session == nil {
		return nil, errors.New("ssh session is nil")
	}

	if len(cmds) == 0 {
		return nil, errors.New("no commands provided")
	}

	cmd := strings.Join(cmds, "; ")
	return session.CombinedOutput(cmd)
}
