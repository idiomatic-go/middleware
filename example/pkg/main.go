package main

import (
	"context"
	"fmt"
	"github.com/idiomatic-go/middleware/example/pkg/host"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"
)

const (
	addr         = "0.0.0.0:8080"
	writeTimeout = time.Second * 60
	readTimeout  = time.Second * 15
	idleTimeout  = time.Second * 60
)

func main() {
	displayRuntime()
	//if !host.Startup() {
	//	os.Exit(1)
	//}
	//defer host.Shutdown()

	r := http.NewServeMux()
	if !host.Startup(r) {
		os.Exit(1)
	}
	
	srv := http.Server{
		Addr: addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      r,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		} else {
			log.Printf("HTTP server Shutdown")
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

func displayRuntime() {
	fmt.Println(fmt.Sprintf("vers : %v", runtime.Version()))
	fmt.Println(fmt.Sprintf("os   : %v", runtime.GOOS))
	fmt.Println(fmt.Sprintf("arch : %v", runtime.GOARCH))
	fmt.Println(fmt.Sprintf("cpu  : %v", runtime.NumCPU()))
}
