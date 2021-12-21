package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/EvisuXiao/andrews-common/config"
	"github.com/EvisuXiao/andrews-common/logging"
	"github.com/EvisuXiao/andrews-common/utils"
)

type IServer interface {
	Config() *config.Server
	Start() error
	Stop() error
	OnStop()
}

type Runner struct {
	srv     IServer
	signal  chan os.Signal
	process chan bool
}

type Option struct {
	Config      *config.Server
	QuitHandler func()
}

func StartServer(s IServer) {
	r := Runner{
		srv:     s,
		signal:  make(chan os.Signal, 1),
		process: make(chan bool, 1),
	}
	r.Run()
}

func (r *Runner) Run() {
	go r.run()
	go r.signalHandler()
	<-r.process
}

func (r *Runner) Stop() {
	logging.Info("Server is quitting")
	ctx, cancel := context.WithTimeout(context.Background(), r.srv.Config().Timeout.Exit)
	defer cancel()
	go r.srv.OnStop()
	if err := r.srv.Stop(); utils.HasErr(err) {
		logging.Error("Stop server err: %+v", err)
		close(r.process)
		return
	}
	select {
	case <-ctx.Done():
		logging.Info("Stop server gracefully")
	}
	close(r.process)
}

func (r *Runner) run() {
	logging.Info("Service %s is running!", config.GetServiceName())
	logging.Info("Start server with listening port %d", config.GetServerConfig().Port)
	logging.Info("The process id is %d", os.Getpid())
	if err := r.srv.Start(); utils.HasErr(err) {
		logging.Error("Start server err: %+v", err)
		close(r.process)
	}
}

func (r *Runner) signalHandler() {
	signal.Notify(r.signal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT)
	select {
	case <-r.signal:
		r.Stop()
	}
}

func NewOption() *Option {
	return &Option{Config: config.GetServerConfig()}
}

func (o *Option) WithConfig(cfg *config.Server) *Option {
	o.Config = cfg
	return o
}

func (o *Option) WithQuitHandler(handler func()) *Option {
	o.QuitHandler = handler
	return o
}
