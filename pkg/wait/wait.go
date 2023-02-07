package wait

import (
	"os"
	"os/signal"
	"syscall"
)

func Interrupt() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}
