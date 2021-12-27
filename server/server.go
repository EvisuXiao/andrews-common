package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	common "github.com/EvisuXiao/andrews-common"
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
	if r.srv.Config().Discovery {
		initDiscoverer()
	}
	go r.run()
	go r.signalHandler()
	<-r.process
}

func (r *Runner) Stop() {
	logging.Info("Server is quitting")
	ctx, cancel := context.WithTimeout(context.Background(), r.srv.Config().Timeout.Exit)
	defer cancel()
	go r.onStop()
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

func (r *Runner) onStop() {
	if err := discoverer.UnregisterInstance(r.srv.Config().Port); utils.HasErr(err) {
		logging.Error("Unregister server instance err: %+v", err)
	}
	common.Stop()
	r.srv.OnStop()
}

func (r *Runner) run() {
	port := r.srv.Config().Port
	weight := r.srv.Config().Weight
	logging.Info("Service %s is running!", config.GetServiceName())
	logging.Info("Start server with listening port %d, weight %f", port, weight)
	logging.Info("The process id is %d", os.Getpid())
	if err := discoverer.RegisterInstance(port, weight, r.buildMetaData()); utils.HasErr(err) {
		logging.Error("Register server instance err: %+v", err)
		close(r.process)
	}
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

func (r *Runner) buildMetaData() map[string]string {
	machine, _ := os.Hostname()
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	return map[string]string{
		"register_time": utils.LocalTimeStr(time.Now()),
		"app_env":       config.GetEnv(),
		"go_version":    runtime.Version(),
		"os":            runtime.GOOS,
		"arch":          runtime.GOARCH,
		"machine_name":  machine,
		"pid":           fmt.Sprint(os.Getpid()),
		"cpu_num":       fmt.Sprint(runtime.NumCPU()),
		"memory":        fmt.Sprint(memStats.Sys),
		"goroutine_num": fmt.Sprint(runtime.NumGoroutine()),
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
