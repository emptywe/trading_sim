package wait

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitInterrupt() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}
