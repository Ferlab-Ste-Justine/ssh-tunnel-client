package ssh

import (
	"golang.org/x/crypto/ssh"
)

func ParsePrivateKey(key []byte) (*ssh.AuthMethod, error) {
	parsedKey, parseKeyErr := ssh.ParsePrivateKey(key)
	if parseKeyErr != nil {
		return nil, parseKeyErr
	}

	auth := ssh.PublicKeys(parsedKey)

	return &auth, nil
}