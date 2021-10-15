package ssh

import (
    "golang.org/x/crypto/ssh"
)

func parsePrivateKey(key []byte) (*ssh.AuthMethod, error) {
    parsedKey, parseKeyErr := ssh.ParsePrivateKey(key)
    if parseKeyErr != nil {
        return nil, parseKeyErr
    }

    auth := ssh.PublicKeys(parsedKey)

    return &auth, nil
}

func processPassword(password []byte) *ssh.AuthMethod {
    authMethod := ssh.Password(string(password))
    return &authMethod
}

func GetAuthMethod(authSecret []byte, authType string) (*ssh.AuthMethod, error) {
    if authType == "key" {
        return parsePrivateKey(authSecret)
    }

    return processPassword(authSecret ), nil
}