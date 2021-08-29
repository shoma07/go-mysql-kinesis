package signal

import (
	"os"
	"os/signal"
	"syscall"
)

type SignalReceiver struct {
	ch chan os.Signal
}

func NewSignalReceiver() (*SignalReceiver, error) {
	sr := &SignalReceiver{}
	ch := make(chan os.Signal, 1)
	signal.Notify(
		ch,
		os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	sr.ch = ch
	return sr, nil
}

func (sr *SignalReceiver) Receive() chan os.Signal {
	return sr.ch
}
