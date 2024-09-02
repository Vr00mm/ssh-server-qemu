package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/phuslu/log"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

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

func (s *SSHServer) handleConnection(nConn net.Conn) {
	defer nConn.Close()

	conn, chans, reqs, err := ssh.NewServerConn(nConn, s.sshConfig)
	if err != nil {
		log.Error().Err(err).Msg("failed to handshake")
		return
	}

	log.Info().Str("remote", conn.RemoteAddr().String()).Str("user", conn.User()).Msg("New SSH connection")

	var wg sync.WaitGroup
	defer wg.Wait()

	// Handle global SSH requests (e.g., keep-alives)
	wg.Add(1)
	go func() {
		ssh.DiscardRequests(reqs)
		wg.Done()
	}()

	// Handle SSH channels
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Error().Err(err).Msg("could not accept channel")
			continue
		}

		term := term.NewTerminal(channel, "> ")

		wg.Add(1)
		go func() {
			defer func() {
				channel.Close()
				wg.Done()
			}()
			for req := range requests {
				if req.Type == "shell" {
					req.Reply(true, nil)
					break
				}
			}

			for {
				line, err := term.ReadLine()
				if err != nil {
					break
				}
				fmt.Println(line) // Process the input (you can customize this part)
			}
		}()
	}
}
