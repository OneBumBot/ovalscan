package ssh

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/onebumbot/ovalscan/internal/config"
	"github.com/onebumbot/ovalscan/internal/utils"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func getKey(keyPath string) ([]byte, error) {

	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, errors.New("unable to read private key")
	}

	return key, nil
}

func createSigner(key []byte, passphrase string) (ssh.Signer, error) {
	signer, err := ssh.ParsePrivateKeyWithPassphrase(key, []byte(passphrase))
	if err != nil {
		// Fallback for unencrypted keys.
		signer, err = ssh.ParsePrivateKey(key)
	}
	if err != nil {
		return nil, errors.New("unable to parse private key:")
	}

	return signer, nil
}

func GetClient(cfg config.SSHConfig) (*ssh.Client, error) {

	keyPath, err := utils.ExpandPath(cfg.KeyPath)
	if err != nil {
		log.Fatal("unable to resolve key path: ", err)
	}

	key, err := getKey(keyPath)

	if err != nil {
		return nil, err
	}

	signer, err := createSigner(key, cfg.Passphrase)

	if err != nil {
		return nil, err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.New("can not find user dir")
	}

	hostKeyCallback, err := knownhosts.New(filepath.Join(homeDir, ".ssh", "known_hosts"))
	if err != nil {
		log.Fatalf("unable to load known_hosts file: %v", err)
	}

	algorithms := ssh.SupportedAlgorithms()
	config := &ssh.ClientConfig{
		Config: ssh.Config{
			KeyExchanges: algorithms.KeyExchanges,
			Ciphers:      algorithms.Ciphers,
			MACs:         algorithms.MACs,
		},
		User: cfg.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback:   hostKeyCallback,
		HostKeyAlgorithms: algorithms.HostKeys,
	}

	client, err := ssh.Dial("tcp", utils.NormalizeHost(cfg.Host), config)
	if err != nil {
		return nil, errors.New("failed to dial")
	}

	return client, nil

}
