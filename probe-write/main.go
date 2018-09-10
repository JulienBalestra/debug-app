package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultListenerPort = 8787
)

type healthCheck struct {
	sleepDelay time.Duration
}

func (h *healthCheck) healthHandler(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("/bin/sh", "-c", "cat /dev/urandom | tr -dc 'a-zA-Z0-9'")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Start()
	// don't wait
	log.Printf("Sleeping %s", h.sleepDelay)
	time.Sleep(h.sleepDelay)
}

func main() {
	execProbe := flag.Bool("health", false, "use this to probe the http listener")
	listenerPort := flag.Int("port", defaultListenerPort, "specify the port for the http listener")
	sleepDelay := flag.Int("sleep", 120, "specify the seconds to sleep over the probe")
	flag.Parse()
	sigCh := make(chan os.Signal)
	defer close(sigCh)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	if *execProbe {
		log.Println("Starting probe")
		c := &http.Client{}
		resp, err := c.Get(fmt.Sprintf("http://127.0.0.1:%d/health", *listenerPort))
		if err != nil {
			log.Fatalf("Error while probing: %v", err)
			// exit 1
		}
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Bad exit code: %d", resp.StatusCode)
			// exit 1
		}
		log.Printf("Successfully probed")
		return
	}

	log.Println("Starting listener ...")
	health := &healthCheck{
		sleepDelay: time.Second * time.Duration(*sleepDelay),
	}
	http.HandleFunc("/health", health.healthHandler)
	listenerBind := fmt.Sprintf("0.0.0.0:%d", *listenerPort)
	log.Printf("Starting to listen on %s", listenerBind)
	go log.Fatalf("%v", http.ListenAndServe(listenerBind, nil))

	for {
		select {
		case s := <-sigCh:
			log.Printf("Ignoring %s", s)
		}
	}
}
