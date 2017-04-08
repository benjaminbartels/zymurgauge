package http

import (
	"net"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var addr = ":3000"

// ToDo: Investigate 1.8 server

// Server is an HTTP server.
type Server struct {
	listener net.Listener
	Handler  *Handler
	Addr     string
	logger   *logrus.Logger
}

// NewServer returns a new instance of the Server
func NewServer(logger *logrus.Logger) *Server {
	return &Server{
		Addr:   addr,
		logger: logger,
	}
}

// Open opens a socket and serves the HTTP server
func (s *Server) Open() error {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return errors.Wrapf(err, "Could not listen on %s", s.Addr)
	}
	s.listener = l
	go func() {
		_ = http.Serve(s.listener, s.Handler)
	}()

	return nil
}

// Close closes the socket
func (s *Server) Close() error {
	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}

// Port returns the port that the server is open on
func (s *Server) Port() int {
	return s.listener.Addr().(*net.TCPAddr).Port
}
