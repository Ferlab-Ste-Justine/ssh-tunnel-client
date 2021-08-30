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

func handleForcedTermination(tunnels []*ssh.SshTunnel) {
    c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(){
		for _ = range c {
            for _, tunnel := range tunnels {
                tunnel.Stop()
            }
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

	ingressTunnel := &ssh.SshTunnel{
		Config: sshConfig,
		LocalFrontendUrl: "127.0.0.1:443",
		SshServerUrl: tunnelServerUrl,
		RemoteBackendUrl: "127.0.0.1:443",
	}
	apiTunnel := &ssh.SshTunnel{
		Config: sshConfig,
		LocalFrontendUrl: "127.0.0.1:6443",
		SshServerUrl: tunnelServerUrl,
		RemoteBackendUrl: "127.0.0.1:6443",
	}

	handleForcedTermination([]*ssh.SshTunnel{ingressTunnel, apiTunnel})

    done := make(chan struct{})

	go func() {
	    err = ingressTunnel.Launch()
	    if err != nil {
		    fmt.Println("Ingress -> Failed to tunnel: " + err.Error())
			if !apiTunnel.IsStopped() {
				apiTunnel.Stop()
			}
	    }
        done<-struct{}{}
	}()

	go func() {
	    err = apiTunnel.Launch()
	    if err != nil {
		    fmt.Println("Api -> Failed to tunnel: " + err.Error())
			if !ingressTunnel.IsStopped() {
				ingressTunnel.Stop()
			}
	    }
        done<-struct{}{}
	}()

    <-done
    <-done
}
