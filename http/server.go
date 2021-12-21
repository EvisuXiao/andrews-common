package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/EvisuXiao/andrews-common/config"
	"github.com/EvisuXiao/andrews-common/server"
	"github.com/EvisuXiao/andrews-common/utils"
)

type Server struct {
	srv    *http.Server
	option *Option
}

type Option struct {
	server.Option
	RouterGroups []*MainRouterGroup
}

func StartServer(option *Option) {
	s := &Server{option: option}
	s.srv = &http.Server{
		Addr:           fmt.Sprintf(":%d", option.Config.Port),
		Handler:        InitRouter(option.RouterGroups...),
		ReadTimeout:    option.Config.Timeout.Read,
		WriteTimeout:   option.Config.Timeout.Write,
		MaxHeaderBytes: 1 << 20,
	}
	server.StartServer(s)
}

func (s *Server) Config() *config.Server {
	return s.option.Config
}

func (s *Server) Start() error {
	if err := s.srv.ListenAndServe(); utils.HasErr(err) && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	if !utils.IsEmpty(s.option.QuitHandler) {
		s.srv.RegisterOnShutdown(s.option.QuitHandler)
	}
	return s.srv.Shutdown(context.Background())
}

func (s *Server) OnStop() {
	// 将handler放在Stop()中进行RegisterOnShutdown()
}

func NewOption(routers ...*MainRouterGroup) *Option {
	o := &Option{Option: *server.NewOption()}
	o.RouterGroups = routers
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
