package runner

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	Interval    int
	StopTimeout int
}

type IRunner interface {
	Run(runAt time.Time)
	StopTimeout()
}

type Runner struct {
	runner IRunner
	conf   *Config
	stopCh chan chan bool
}

func New(conf *Config, runner IRunner) *Runner {
	return &Runner{
		runner: runner,
		conf:   conf,
		stopCh: make(chan chan bool, 1),
	}
}

func (self *Runner) Start() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	go self.startTicker()
	<-sigCh

	stopTimeout := time.Duration(self.conf.StopTimeout) * time.Second
	stopCh := make(chan bool, 1)
	go func() {
		self.stop()
		stopCh <- true
	}()
	select {
	case <-stopCh:
	case <-time.After(stopTimeout):
		self.runner.StopTimeout()
	}
}

func (self *Runner) startTicker() {
	interval := time.Duration(self.conf.Interval) * time.Second
	ticker := time.NewTicker(interval)
	for {
		select {
		case doneCh := <-self.stopCh:
			doneCh <- true
			return
		case t := <-ticker.C:
			self.runner.Run(t)
		}
	}
}

func (self *Runner) stop() {
	doneCh := make(chan bool)
	self.stopCh <- doneCh
	<-doneCh
}
