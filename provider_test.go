package main

import (
	"testing"
)

func TestProvider(t *testing.T) {
	p := Provider()
	if err := p.InternalValidate(); err != nil {
		t.Fatalf("Provider internal validation failed: %s", err)
	}
}
