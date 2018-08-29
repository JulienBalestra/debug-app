package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"sync"
	"syscall"
	"time"
)

const (
	defaultHealthDir    = "/var/lib/debug-app/"
	defaultListenerPort = 8787
)

type healthCheck struct {
	healthDir string
	sync.Mutex
	healthy     bool
	maxFork     int
	currentFork int
	probeSleep  bool
}

func (h *healthCheck) healthHandler(w http.ResponseWriter, r *http.Request) {
	h.Lock()
	defer h.Unlock()
	if !h.healthy {
		log.Printf("Marked as unhealthy, returns 500")
		w.WriteHeader(500)
		return
	}

	// read correctly
	healthFile := path.Join(h.healthDir, "health-file")
	b, err := ioutil.ReadFile(healthFile)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Fail to read %s: %v", healthFile, err)
		w.WriteHeader(500)
		return
	}
	log.Printf("health check file read: %q", string(b))

	// write correctly
	healthNow := time.Now().Format("2006-01-02T15:04:05Z")
	healthNowB := []byte(healthNow)
	err = ioutil.WriteFile(healthFile, healthNowB, os.ModePerm)
	if err != nil {
		log.Printf("Fail to write %s: %v", healthFile, err)
		w.WriteHeader(500)
		return
	}
	log.Printf("health check file wrote: %q", healthNow)

	// jitter timeouts
	randInt := rand.Intn(10)
	if h.probeSleep {
		sec := time.Duration(randInt) * time.Second
		log.Printf("Sleeping %s", sec.String())
		<-time.After(sec)
	}

	// dirty open: never close fd
	f, err := os.OpenFile(fmt.Sprintf("/tmp/dirty-%d", randInt), os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Printf("Unexpected error while dirty open: %v", err)
		w.WriteHeader(500)
		return
	}
	_, err = f.Write(healthNowB)
	if err != nil {
		log.Printf("Fail to write in dirty open %v", err)
		w.WriteHeader(500)
		return
	}
	log.Println("Health check returns 200")
	// don't flush, don't close: this is dirty
}

func (h *healthCheck) switchHealthyHandler(_ http.ResponseWriter, _ *http.Request) {
	h.Lock()
	defer h.Unlock()
	log.Printf("Healthy: %v", h.healthy)
	h.healthy = !h.healthy
	log.Printf("Switched as healthy: %v", h.healthy)
}

func (h *healthCheck) forkHandler(w http.ResponseWriter, _ *http.Request) {
	h.Lock()
	defer h.Unlock()
	if h.currentFork >= h.maxFork {
		w.WriteHeader(403)
		return
	}
	h.currentFork++
	log.Printf("Fork OK: %d/%d", h.currentFork, h.maxFork)
}

func main() {
	execProbe := flag.Bool("health", false, "use this to probe the http listener")
	probeSleep := flag.Bool("probe-sleep", false, "use this to introduce a random sleep during probe")
	listenerPort := flag.Int("port", defaultListenerPort, "specify the port for the http listener")
	healthDir := flag.String("health-dir", defaultHealthDir, "specify the dir for the health handler")
	maxFork := flag.Int("fork", 3, "specify the number of dirty fork during probing")
	flag.Parse()

	if *execProbe {
		log.Println("Starting probe")
		c := &http.Client{
			Timeout: time.Second * 30,
		}
		// fork but don't wait
		resp, err := c.Get(fmt.Sprintf("http://127.0.0.1:%d/fork", *listenerPort))
		if err != nil {
			log.Fatalf("Error while probing: %v", err)
			// exit 1
		}
		if resp.StatusCode == http.StatusOK {
			exec.Command("/bin/sh", "-c", "nohup tail -f "+path.Join(*healthDir, "health-file")).Start()
		}

		resp, err = c.Get(fmt.Sprintf("http://127.0.0.1:%d/health", *listenerPort))
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
	log.Printf("Creating %s ..", *healthDir)
	err := os.MkdirAll(*healthDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("Cannot create dir: %v", err)
		// exit 1
	}
	sigCh := make(chan os.Signal, 2)
	defer close(sigCh)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	health := &healthCheck{
		healthDir:  *healthDir,
		healthy:    true,
		maxFork:    *maxFork,
		probeSleep: *probeSleep,
	}
	http.HandleFunc("/health", health.healthHandler)
	http.HandleFunc("/fork", health.forkHandler)
	http.HandleFunc("/switch", health.switchHealthyHandler)
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
