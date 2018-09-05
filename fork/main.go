package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
	"os/exec"
	"log"
)

func main() {
	sigCh := make(chan os.Signal, 2)
	defer close(sigCh)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()
	for {
		select {
		case <-sigCh:
			return
		case <-ticker.C:
			b, err := exec.Command("/bin/sh", "-c", "cat /etc/hosts | cat -e").CombinedOutput()
			if err != nil {
				log.Printf("Error: %v\n", err)
			}
			log.Printf("\n%s", string(b))
		}
	}
}
