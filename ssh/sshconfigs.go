package ssh

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"

	"golang.org/x/crypto/ssh"
)

func GetAuthSshConfigs() (*ssh.ClientConfig, error) {
	buffer, err := ioutil.ReadFile("authorized-ssh-private-key")
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}

	buffer, err = ioutil.ReadFile("host-md5-fingerprint")
	if err != nil {
		return nil, err
	}
	md5Fingerprint := string(buffer)
	
	return &ssh.ClientConfig{
		User: "ubuntu",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			//rework later to use ssh.FingerprintSHA256(key)
			hostFingerprint := ssh.FingerprintLegacyMD5(key)
			if hostFingerprint != md5Fingerprint {
				return errors.New(fmt.Sprintln("Server %s fingerprint did not match expected %s fingerprint", hostFingerprint, md5Fingerprint))
			}
			return nil
		},
	}, nil
}