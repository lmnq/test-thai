package fasthttpserver

import (
	"time"

	"github.com/valyala/fasthttp"
)

const (
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultIdleTimeout     = 10 * time.Second
	_defaultAddr            = ":8080"
	_defaultShutdownTimeout = 3 * time.Second
)

// Server -.
type Server struct {
	server *fasthttp.Server
	notify chan error
}

// NewServer -.
func New(handler fasthttp.RequestHandler, port string, opts ...Option) *Server {
	fasthttpServer := &fasthttp.Server{
		Handler:      handler,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
		IdleTimeout:  _defaultIdleTimeout,
	}

	s := &Server{
		server: fasthttpServer,
		notify: make(chan error, 1),
	}

	// custom options
	for _, opt := range opts {
		opt(s)
	}

	s.start(port)

	return s
}

func (s *Server) start(port string) {
	if port == "" {
		port = _defaultAddr
	}

	go func() {
		err := s.server.ListenAndServe(":" + port)
		s.notify <- err
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	return s.server.Shutdown()
}
