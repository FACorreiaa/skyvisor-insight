package config

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	"time"
)

func InitPprof(addr, port string) error {
	if addr == "" || port == "" {
		return nil
	}
	ip := net.ParseIP(addr)
	if ip == nil || !ip.IsLoopback() {
		return errors.New("pprof must bind to a loopback address")
	}
	pprofAddr := fmt.Sprintf("%s:%s", addr, port)
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	server := &http.Server{
		Addr:         pprofAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
		Handler:      mux,
	}
	go func() {
		slog.Debug("pprof server listening", "address", pprofAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("pprof server failed", "error", err)
		}
	}()
	return nil
}
