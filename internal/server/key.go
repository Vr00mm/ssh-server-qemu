package server

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
	"github.com/phuslu/log"
)

// GeneratePrivateKey generates an RSA private key.
func GeneratePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %v", err)
	}
	return privateKey, nil
}

// GenerateSSHSigner generates an SSH signer from an RSA private key.
func GenerateSSHSigner(privateKey *rsa.PrivateKey) (ssh.Signer, error) {
	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer: %v", err)
	}
	return signer, nil
}

// ParseSSHPrivateKeyFromBytes parses an SSH private key from a byte slice.
func ParseSSHPrivateKeyFromBytes(privateKeyBytes []byte) (ssh.Signer, error) {
	private, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}
	return private, nil
}

// LoadOrGeneratePrivateKey tries to load the private key from a file.
// If it fails, it generates a new one and returns it without writing to disk.
func (s *SSHServer) LoadOrGeneratePrivateKey() (ssh.Signer, error) {
	privateBytes, err := os.ReadFile(s.config.HostPrivKey)
	if err != nil {
		log.Warn().Msg("Private key not found or unreadable, generating a new one")

		privateKey, err := GeneratePrivateKey(2048)
		if err != nil {
			return nil, err
		}

		signer, err := GenerateSSHSigner(privateKey)
		if err != nil {
			return nil, err
		}

		log.Info().Msg("Generated new private key in memory")
		return signer, nil
	}

	return ParseSSHPrivateKeyFromBytes(privateBytes)
}