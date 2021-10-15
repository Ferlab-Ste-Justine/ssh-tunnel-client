package main

import (
	"fmt"
	_ "embed"
	"io/ioutil"
	"os"
	"os/signal"
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

func getAuthSecret(fallbackVal string) ([]byte, error) {
	if _, err := os.Stat("auth_secret"); err != nil && fallbackVal != "" {
		return []byte(fallbackVal), nil
	}

	buffer, err := ioutil.ReadFile("auth_secret")
	if err != nil {
		return []byte{}, err
	}
	return buffer, nil
}

var (
    //go:embed auth_secret
    authSecret string
	//go:embed tunnel_config.json
	tunnelConfig string
)

func main() {
	tunnelConfig, tunnelConfigErr := getTunnelConfig(tunnelConfig)
	if tunnelConfigErr != nil {
		panic("Failed to get tunnel configuration: " + tunnelConfigErr.Error())
	}

    authSecretAsBytes, authSecretAsBytesErr := getAuthSecret(authSecret)
	if authSecretAsBytesErr != nil {
		panic("Failed to read auth secret: " + authSecretAsBytesErr.Error())
	}

	authMethod, authMethodErr := ssh.GetAuthMethod(authSecretAsBytes, tunnelConfig.AuthMethod)
	if authMethodErr != nil {
		panic("Failed to parse authentication secret: " + authMethodErr.Error())
	}

	sshConfig, err := ssh.GetAuthSshConfigs(tunnelConfig.HostUser, *authMethod, tunnelConfig.HostMd5FingerPrint)
	if err != nil {
		panic("Failed to initiate ssh auth configs: " + err.Error())
	}

	tunnels := []*ssh.SshTunnel{}
    for _, binding := range tunnelConfig.Bindings {
		tunnels = append(
			tunnels, 
			&ssh.SshTunnel{
			    Config: sshConfig,
			    LocalFrontendUrl: binding.Local,
			    SshServerUrl: tunnelConfig.HostUrl,
			    RemoteBackendUrl: binding.Remote,
			},
	   )
	}
	manager := &ssh.SshTunnelsManager{Tunnels: tunnels}
	handleForcedTermination(manager)

	errs := manager.Launch()
	for _, err := range errs {
		fmt.Println(err.Error())
	}
}
