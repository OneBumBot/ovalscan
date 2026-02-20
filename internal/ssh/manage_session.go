package ssh

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	err := session.RequestPty("xterm", 80, 40, modes)

	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("request pty failed: %w (session may already be used or closed; create a new session)", err)
		}
		return nil, err
	}

	in, err := session.StdinPipe()

	if err != nil {
		return nil, err
	}

	out, err := session.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := session.Shell(); err != nil {
		return nil, err
	}

	cmd := strings.Join(cmds, "; ")
	if _, err := io.WriteString(in, cmd+"\nexit\n"); err != nil {
		return nil, err
	}

	if err := in.Close(); err != nil {
		return nil, err
	}

	var output []byte
	r := bufio.NewReader(out)
	for {
		b, readErr := r.ReadByte()
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				break
			}
			return nil, readErr
		}
		output = append(output, b)
	}

	if err := session.Wait(); err != nil {
		return output, err
	}

	return output, nil
}
