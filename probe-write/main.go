package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	defaultListenerPort = 8787
)

type healthCheck struct{}

func (h *healthCheck) healthHandler(w http.ResponseWriter, r *http.Request) {
	urandom, err := os.Open("/dev/urandom")
	if err != nil {
		panic(err)
	}
	defer urandom.Close()
	io.Copy(os.Stdout, urandom)
}

func main() {
	execProbe := flag.Bool("health", false, "use this to probe the http listener")
	listenerPort := flag.Int("port", defaultListenerPort, "specify the port for the http listener")
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
	health := &healthCheck{}
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
