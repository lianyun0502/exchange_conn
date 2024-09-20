package common

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func WaitForClose(log *logrus.Logger, stopCh chan struct{}) {
	sysSignalCh := make(chan os.Signal, 2)
	signal.Notify(sysSignalCh, os.Interrupt, syscall.SIGTERM)
	select {
	case sig := <-sysSignalCh:
		log.Info("receive signal: ", sig)
		close(stopCh)
	case <- stopCh:
		log.Info("receive close signal")
	}
	log.Info("Closed")
}