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

	for {
		frontendConn, err := tunnel.frontendListener.Accept()
		if err != nil {
			if !tunnel.stopped {
				return err
			}
			return nil
		}
		go tunnel.forwardFrontendToBackend(frontendConn)
	}
}

func pipeConn(writer net.Conn, reader net.Conn, c chan error) {
	_, err := io.Copy(writer, reader)
	c <- err
}


func (tunnel *SshTunnel) forwardFrontendToBackend(frontendConn net.Conn) {
	defer frontendConn.Close()

	sshConn, err := ssh.Dial("tcp", tunnel.SshServerUrl, tunnel.Config)
	if err != nil {
		fmt.Printf("Ssh dial error: %s\n", err)
		return
	}
	defer sshConn.Close()

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
		fmt.Printf("io.Copy error: %s", err)
	}
	err = <- pipeErrChan
	if err != nil {
		fmt.Printf("io.Copy error: %s", err)
	}
}