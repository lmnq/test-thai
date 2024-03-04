package fasthttpserver

import "time"

// Option -.
type Option func(*Server)

// ReadTimeout -.
func ReadTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = t
	}
}

// WriteTimeout -.
func WriteTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = t
	}
}

// IdleTimeout -.
func IdleTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.server.IdleTimeout = t
	}
}
