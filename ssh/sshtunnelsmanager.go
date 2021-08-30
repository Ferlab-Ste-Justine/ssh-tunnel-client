package ssh

import (
	"errors"
	"fmt"
)


type SshTunnelsManager struct {
    Tunnels []*SshTunnel
	done bool
}

func (manager *SshTunnelsManager) Close() {
	for _, tunnel := range manager.Tunnels {
		if !tunnel.IsClosed() {
			tunnel.Close()
		}
	}
}

func (manager *SshTunnelsManager) Launch() []error {
	var err error

	defer manager.Close()
	for _, tunnel := range manager.Tunnels {
		err = tunnel.Init()
		if err != nil {
			return []error{errors.New(fmt.Sprintf("Tunnel init err: %s", err.Error()))}
		}
	}

    done := make(chan error)
	for _, tunnel := range manager.Tunnels {
		go func(t *SshTunnel) {
			listenErr := t.Listen()
			if listenErr != nil {
                manager.Close()
			}
			done<-listenErr
		}(tunnel)
	}

	errs := []error{}
	for _, _ = range manager.Tunnels {
	    err = <-done
		if err != nil {
			err = errors.New(fmt.Sprintf("Tunnel listen err: %s", err.Error()))
			errs = append(errs, err)
		}
	}

	return errs
}