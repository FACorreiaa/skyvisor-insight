package config

import (
	"fmt"
	"log"
	"net/http"
	"time"

	// nolint:gosec,G108
	_ "net/http/pprof"
)

func InitPprof(addr, port string) error {
	pprofAddr := fmt.Sprintf("%s:%s", addr, port)
	server := &http.Server{
		Addr:         pprofAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		log.Printf("Running pprof on %s\n", pprofAddr)
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Error running pprof server: %v\n", err)
		}
	}()
	return nil
}
