package ssh

import (
	"errors"
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
)

func GetAuthSshConfigs(sshUser string, key ssh.AuthMethod, sha256Fingerprint string) (*ssh.ClientConfig, error) {	
	return &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			key,
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			//rework later to use ssh.FingerprintSHA256(key)
			hostFingerprint := ssh.FingerprintSHA256(key)
			if hostFingerprint != sha256Fingerprint {
				return errors.New(fmt.Sprintln("Server %s fingerprint did not match expected %s fingerprint", hostFingerprint, sha256Fingerprint))
			}
			return nil
		},
	}, nil
}