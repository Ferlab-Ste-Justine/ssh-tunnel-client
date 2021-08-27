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

func main() {
	sshConfig, err := ssh.GetAuthSshConfigs()
	if err != nil {
		panic("Failed to initiate ssh auth configs: " + err.Error())
	}

	tunnelServerUrl, tunnelUrlErr := getTunnelServerUrl()
	if err != nil {
		panic("Failed to retrieve tunnel server url: " + tunnelUrlErr.Error())
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
