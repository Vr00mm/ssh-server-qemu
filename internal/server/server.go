package server

import (
	"fmt"

	"golang.org/x/crypto/ssh"

	"yourusername/sshserver/internal/config"
)

type SSHServer struct {
	config    *config.Config
	sshConfig *ssh.ServerConfig
}

func New(cfg *config.Config) (*SSHServer, error) {
	server := &SSHServer{
		config: cfg,
	}

	if err := server.setupSSHConfig(); err != nil {
		return nil, err
	}

	return server, nil
}

func (s *SSHServer) setupSSHConfig() error {
	signer, err := s.LoadOrGeneratePrivateKey()
	if err != nil {
		return fmt.Errorf("failed to load or generate private key: %w", err)
	}

	s.sshConfig = &ssh.ServerConfig{
		MaxAuthTries:                1,
		NoClientAuth:                false,
		PublicKeyCallback:           nil,
		KeyboardInteractiveCallback: nil,
		PasswordCallback:            s.authenticateUser,
	}

	s.sshConfig.AddHostKey(signer)

	return nil
}
