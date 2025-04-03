package service

import (
	"os"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
)

type Service struct {
	appName         string
	diagnosticsAddr string

	ready     bool
	readyLock sync.RWMutex
}

func (s *Service) SetReady(state bool) {
	s.readyLock.Lock()
	defer s.readyLock.Unlock()
	s.ready = state
}

func (s *Service) IsReady() bool {
	s.readyLock.RLock()
	defer s.readyLock.RUnlock()

	return s.ready
}

type Option func(*Service)

func WithDiagnosticsServer(addr string) Option {
	return func(s *Service) {
		s.diagnosticsAddr = addr
	}
}

func New(appName string, opts ...Option) *Service {
	svc := Service{appName: appName}

	for _, opt := range opts {
		opt(&svc)
	}

	logrus.Info("initializing app")

	return &svc
}

type StartStopper interface {
	Start()
	Stop()
}

func (s *Service) RunWait(services ...StartStopper) *sync.WaitGroup {
	logrus.Info("starting app")

	for _, s := range services {
		s.Start()
	}

	s.SetReady(true)
	logrus.Info("app ready")

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		logrus.WithField("signal", Wait([]os.Signal{syscall.SIGTERM, syscall.SIGINT}).String()).Info("received signal")
		logrus.Info("stopping app")

		s.SetReady(false)

		// stop in reverse order
		for i := range services {
			services[len(services)-i-1].Stop()
		}

		logrus.Info("bye 👋")
	}()

	return wg
}
