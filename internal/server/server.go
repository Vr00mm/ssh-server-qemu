package server

import (
	"fmt"
	"net"

	"github.com/phuslu/log"
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

func (s *SSHServer) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.config.ListenAddress, s.config.SSHPort))
	if err != nil {
		return fmt.Errorf("failed to listen on %s:%d: %w", s.config.ListenAddress, s.config.SSHPort, err)
	}
	defer listener.Close()

	log.Info().Str("address", listener.Addr().String()).Msg("SSH server listening")

	for {
		nConn, err := listener.Accept()
		if err != nil {
			log.Error().Err(err).Msg("failed to accept incoming connection")
			continue
		}
		go s.handleConnection(nConn)
	}
}
