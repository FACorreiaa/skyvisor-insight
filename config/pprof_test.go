package config

import "testing"

func TestInitPprofRejectsPublicAddress(t *testing.T) {
	t.Parallel()

	if err := InitPprof("0.0.0.0", "6060"); err == nil {
		t.Fatal("InitPprof() error = nil, want public-address error")
	}
}

func TestInitPprofDisabledWithoutAddress(t *testing.T) {
	t.Parallel()

	if err := InitPprof("", ""); err != nil {
		t.Fatalf("InitPprof() error = %v", err)
	}
}
