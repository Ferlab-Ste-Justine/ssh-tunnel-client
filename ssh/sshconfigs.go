package ssh

import (
	"errors"
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
)

func GetAuthSshConfigs(sshUser string, key ssh.AuthMethod, md5Fingerprint string) (*ssh.ClientConfig, error) {	
	return &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			key,
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