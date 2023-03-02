package wait

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Interrupt() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	time.Sleep(time.Second * 5)
	<-quit
}
