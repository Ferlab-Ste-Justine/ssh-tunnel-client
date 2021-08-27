package ssh

import (
	"fmt"
	"io"
	"net"

	"golang.org/x/crypto/ssh"
)

type SshTunnel struct {
    LocalFrontendUrl string
    SshServerUrl string
	RemoteBackendUrl string
	Config *ssh.ClientConfig
	frontendListener net.Listener
	stopped bool
}

func (tunnel *SshTunnel) IsStopped() bool {
	return tunnel.stopped
}

func (tunnel *SshTunnel) Stop() {
	tunnel.stopped = true
	tunnel.frontendListener.Close()
}

func (tunnel *SshTunnel) Launch() error {
	var err error
	tunnel.stopped = false
	tunnel.frontendListener, err = net.Listen("tcp", tunnel.LocalFrontendUrl)
	if err != nil {
		return err
	}
	defer tunnel.frontendListener.Close()

	sshConn, sshErr := ssh.Dial("tcp", tunnel.SshServerUrl, tunnel.Config)
	if sshErr != nil {
		return sshErr
	}
	defer sshConn.Close()

	for {
		frontendConn, err := tunnel.frontendListener.Accept()
		if err != nil {
			if !tunnel.stopped {
				return err
			}
			return nil
		}
		go tunnel.forwardFrontendToBackend(frontendConn, sshConn)
	}
}

func pipeConn(writer net.Conn, reader net.Conn, c chan error) {
	_, err := io.Copy(writer, reader)
	c <- err
}


func (tunnel *SshTunnel) forwardFrontendToBackend(frontendConn net.Conn, sshConn *ssh.Client) {
	defer frontendConn.Close()

	backendConn, err := sshConn.Dial("tcp", tunnel.RemoteBackendUrl)
	if err != nil {
		fmt.Printf("Backend dial error: %s\n", err)
		return
	}
	defer backendConn.Close()

	pipeErrChan := make(chan error)

	go pipeConn(frontendConn, backendConn, pipeErrChan)
	go pipeConn(backendConn, frontendConn, pipeErrChan)
	err = <- pipeErrChan
	if err != nil {
		fmt.Printf("io.Copy error: %s\n", err)
	}
	err = <- pipeErrChan
	if err != nil {
		fmt.Printf("io.Copy error: %s\n", err)
	}
}