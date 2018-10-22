package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"math/rand"
	"path/filepath"
	"time"
	"syscall"
	"strconv"
)

const minPort = 1024
const maxPort = 65535

func main() {
	var port int
	var dir string
	var err error
	if dir, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		log.Fatalf("Could not access current directory: %v", err)
	}
	if len(os.Args) > 1 {
		var p int64
		if p, err = strconv.ParseInt(os.Args[1], 10, 64); err != nil {
			log.Fatalf("Invalid port: %s", os.Args[1])
		} else if p < minPort || p > maxPort {
			log.Fatalf("Invalid port: %d, must be [%d - %d]", p, minPort, maxPort)
		} else {
			port = int(p)
		}
	} else {
		rand.Seed(time.Now().UnixNano())
		port = rand.Intn(maxPort - minPort) + minPort
	}
	web := http.Server{Addr: fmt.Sprintf(":%d", port), Handler: http.DefaultServeMux}
	fs := http.FileServer(http.Dir(dir))
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		if err = web.Shutdown(context.Background()); err != nil {
			log.Fatalf("Could not shutdown: %v", err)
		}
	}()

	http.Handle("/", fs)

	log.Printf("Starting on port: %d\n", port)

	if err = web.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Could not start http listener: %v\n", err)
	}
}
