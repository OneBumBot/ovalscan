package ssh

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"log"

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
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	err := session.RequestPty("xterm", 80, 40, modes)

	if err != nil {
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

	log.Print("Pipes successfully created")

	var output []byte

	go func(in io.WriteCloser, out io.Reader, output *[]byte) {
		var (
			line string
			r    = bufio.NewReader(out)
		)

		for {
			b, err := r.ReadByte()
			if err != nil {
				break
			}

			*output = append(*output, b)

			if b == byte('\n') {
				line = ""
				continue
			}

			line += string(b)

			if strings.HasPrefix(line, "[sudo] password for ") && strings.HasSuffix(line, ": ") {
				_, err = in.Write([]byte("123\n"))
				if err != nil {
					break
				}
			}

		}

	}(in, out, &output)

	log.Print("preparing to execute commandss")

	cmd := strings.Join(cmds, "; ")
	_, err = session.Output(cmd)
	if err != nil {
		return nil, err
	}

	return output, nil

}
