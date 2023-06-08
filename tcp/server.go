package tcp

import (
	"context"
	"fmt"
	"ledis/interface/tcp"
	"ledis/lib/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	Address string
}

func ListenAndServerWithSignal(cfg *Config, handler tcp.Handler) error {
	closeChan := make(chan struct{})
	sigChan := make(chan os.Signal)

	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		syscall.SIGINT)
	go func() {
		sig := <-sigChan
		switch sig {
		case syscall.SIGHUP,
			syscall.SIGQUIT,
			syscall.SIGTERM,
			syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()

	listen, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("bind: %s, start listening...", cfg.Address))

	ListenAndServer(listen, handler, closeChan)

	return nil
}

func ListenAndServer(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {

	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()

	go func() {
		<-closeChan
		logger.Info("shutting down")
		_ = listener.Close()
		_ = handler.Close()
	}()

	ctx := context.Background()

	var waitDone sync.WaitGroup
	for true {
		conn, err := listener.Accept()
		if err != nil {
			logger.Info(err)
			break
		}
		logger.Info("accepted new conn")
		waitDone.Add(1)

		go func() {
			defer func() {
				waitDone.Done()
			}()
			handler.Handler(ctx, conn)
		}()
	}

	waitDone.Wait()
}
