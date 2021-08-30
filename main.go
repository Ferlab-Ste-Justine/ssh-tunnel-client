package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"ferlab/k8tunnel/ssh"
)

func handleForcedTermination(manager *ssh.SshTunnelsManager) {
    c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(){
		for _ = range c {
            manager.Close()
		}
	}()
}

func getTunnelServerUrl() (string, error) {
	buffer, err := ioutil.ReadFile("tunnel-server-url")
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(buffer), "\n"), nil
}

func getTunnelServerFingerprint() (string, error) {
	buffer, err := ioutil.ReadFile("host-md5-fingerprint")
	if err != nil {
		return "", err
	}
	return string(buffer), nil
}

func getSshPrivateKey() ([]byte, error) {
	buffer, err := ioutil.ReadFile("authorized-ssh-private-key")
	if err != nil {
		return []byte{}, err
	}
	return buffer, nil
}

func main() {
    privateKeyAsBytes, keyAsBytesErr := getSshPrivateKey()
	if keyAsBytesErr != nil {
		panic("Failed to read private key: " + keyAsBytesErr.Error())
	}

	privateKeyPtr, keyParseErr := ssh.ParsePrivateKey(privateKeyAsBytes)
	if keyParseErr != nil {
		panic("Failed to parse private key: " + keyParseErr.Error())
	}

	tunnelServerUrl, tunnelUrlErr := getTunnelServerUrl()
	if tunnelUrlErr != nil {
		panic("Failed to retrieve tunnel server url: " + tunnelUrlErr.Error())
	}

	tunnelFingerprint, tunnelFingerprintErr := getTunnelServerFingerprint()
	if tunnelFingerprintErr != nil {
		panic("Failed to retrieve tunnel server fingerprint: " + tunnelFingerprintErr.Error())
	}

	sshConfig, err := ssh.GetAuthSshConfigs("ubuntu", *privateKeyPtr, tunnelFingerprint)
	if err != nil {
		panic("Failed to initiate ssh auth configs: " + err.Error())
	}

	tunnels := []*ssh.SshTunnel{
		&ssh.SshTunnel{
			Config: sshConfig,
			LocalFrontendUrl: "127.0.0.1:443",
			SshServerUrl: tunnelServerUrl,
			RemoteBackendUrl: "127.0.0.1:443",
		},
		&ssh.SshTunnel{
			Config: sshConfig,
			LocalFrontendUrl: "127.0.0.1:6443",
			SshServerUrl: tunnelServerUrl,
			RemoteBackendUrl: "127.0.0.1:6443",
		},
	}
	manager := &ssh.SshTunnelsManager{Tunnels: tunnels}
	handleForcedTermination(manager)

	errs := manager.Launch()
	for _, err := range errs {
		fmt.Println(err.Error())
	}
}
