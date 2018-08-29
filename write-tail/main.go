package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func startWriting() error {
	// dirty open: never close fd
	f, err := os.OpenFile("/tmp/write-tail", os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Printf("Unexpected error while open: %v", err)
		return err
	}

	go func() {
		ticker := time.NewTicker(time.Millisecond * 100)
		defer ticker.Stop()
		log.Printf("Start writing")
		for {
			select {
			case <-ticker.C:
				healthNow := time.Now().Format("2006-01-02T15:04:05Z")
				healthNowB := []byte(healthNow)
				_, err = f.Write(healthNowB)
				if err != nil {
					log.Printf("Fail to write in dirty open %v", err)
				}
			}
		}
	}()
	return nil
}

func main() {
	err := startWriting()
	if err != nil {
		log.Fatalf("Cannot start writing: %v", err)
		// exit 1
	}

	// fork but don't wait read / write
	exec.Command("/bin/sh", "-c", "nohup tail -f /tmp/write-tail").Start()
	_, err = os.Create("/tmp/d")
	if err != nil {
		log.Fatalf("Cannot create file: %v", err)
		// exit 1
	}
	exec.Command("/bin/sh", "-c", "exec nohup sh -c 'while true; do date > /tmp/d; done'").Start()
	exec.Command("/bin/sh", "-c", "nohup tail -f /tmp/d").Start()

	sigCh := make(chan os.Signal, 2)
	defer close(sigCh)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case s := <-sigCh:
			log.Printf("Ignoring %s", s)
		}
	}
}
