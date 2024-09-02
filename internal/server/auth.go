package server

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/phuslu/log"
	"golang.org/x/crypto/ssh"
)

func (s *SSHServer) authenticateUser(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
	log.Info().Str("user", conn.User()).Msg("Authenticating user")

	resp, err := http.PostForm(s.config.AuthnURL, url.Values{
		"username": {conn.User()},
		"password": {string(password)},
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to contact authentication server")
		return nil, fmt.Errorf("failed to contact authentication server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn().Str("user", conn.User()).Int("status", resp.StatusCode).Msg("Authentication failed")
		return nil, fmt.Errorf("authentication failed with status code: %d", resp.StatusCode)
	}

	log.Info().Str("user", conn.User()).Msg("Authentication successful")
	return nil, nil
}
