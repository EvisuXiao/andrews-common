package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/EvisuXiao/andrews-common/config"
	"github.com/EvisuXiao/andrews-common/logging"
	"github.com/EvisuXiao/andrews-common/server"
	"github.com/EvisuXiao/andrews-common/utils"
)

type Server struct {
	srv      *grpc.Server
	listener net.Listener
	option   *Option
}

type Option struct {
	server.Option
}

func StartServer(option *Option) {
	port := fmt.Sprintf(":%d", option.Config.Port)
	listener, err := net.Listen("tcp", port)
	if utils.HasErr(err) {
		logging.Fatal("Listen tcp port(%d) err: %+v", port, err)
	}
	s := &Server{srv: grpc.NewServer(), listener: listener, option: option}
	s.srv = grpc.NewServer()
	server.StartServer(s)
}

func (s *Server) Config() *config.Server {
	return s.option.Config
}

func (s *Server) Start() error {
	return s.srv.Serve(s.listener)
}

func (s *Server) Stop() error {
	s.srv.GracefulStop()
	return nil
}

func (s *Server) OnStop() {
	if !utils.IsEmpty(s.option.QuitHandler) {
		go s.option.QuitHandler()
	}
}

func NewOption() *Option {
	o := &Option{Option: *server.NewOption()}
	return o
}

func (o *Option) WithConfig(cfg *config.Server) *Option {
	o.Option.WithConfig(cfg)
	return o
}

func (o *Option) WithQuitHandler(handler func()) *Option {
	o.Option.WithQuitHandler(handler)
	return o
}
