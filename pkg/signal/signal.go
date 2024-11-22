package signal

import (
	"os"
	"os/signal"
	"syscall"
)

func Func(callback func()) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	callback()
}
