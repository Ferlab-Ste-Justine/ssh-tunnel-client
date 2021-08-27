package ssh

import (
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

func ReadPrivateKeyFromFile(path string) (*ssh.AuthMethod, error) {
	buffer, readFileErr := ioutil.ReadFile(path)
	if readFileErr != nil {
		return nil, readFileErr
	}

	key, parseKeyErr := ssh.ParsePrivateKey(buffer)
	if parseKeyErr != nil {
		return nil, parseKeyErr
	}

	auth := ssh.PublicKeys(key)

	return &auth, nil
}