package ssh

import (
	"errors"
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
	sshConn *ssh.Client
	listening bool
}

func (tunnel *SshTunnel) IsClosed() bool {
	return !tunnel.listening
}

func (tunnel *SshTunnel) Close() {
	tunnel.listening = false
	tunnel.frontendListener.Close()
	tunnel.sshConn.Close()
}

func (tunnel *SshTunnel) Init() error {
    var err error

	tunnel.frontendListener, err = net.Listen("tcp", tunnel.LocalFrontendUrl)
	if err != nil {
		return err
	}

	tunnel.sshConn, err = ssh.Dial("tcp", tunnel.SshServerUrl, tunnel.Config)
	if err != nil {
		tunnel.frontendListener.Close()
		return err
	}
	tunnel.listening = true

	return nil
}

func (tunnel *SshTunnel) Listen() error {
	for {
		frontendConn, err := tunnel.frontendListener.Accept()
		if err != nil {
			if tunnel.listening {
				return err
			}
			return nil
		}
		go tunnel.forwardFrontendToBackend(frontendConn)
	}
}

func pipeConn(writer net.Conn, reader net.Conn, c chan error) {
	defer writer.Close()
	defer reader.Close()

	_, err := io.Copy(writer, reader)
	c <- err
}


func (tunnel *SshTunnel) forwardFrontendToBackend(frontendConn net.Conn) {
	backendConn, err := tunnel.sshConn.Dial("tcp", tunnel.RemoteBackendUrl)
	if err != nil {
		fmt.Printf("Backend dial error: %s\n", err)
		return
	}

	pipeErrChan := make(chan error)

	go pipeConn(frontendConn, backendConn, pipeErrChan)
	go pipeConn(backendConn, frontendConn, pipeErrChan)
	err = <- pipeErrChan
	if err != nil && !errors.Is(err, net.ErrClosed) {
		fmt.Printf("%s\n", err)
	}
	err = <- pipeErrChan
	if err != nil && !errors.Is(err, net.ErrClosed) {
		fmt.Printf("%s\n", err)
	}
}