package main

import "testing"

func TestAlwaysFails(t *testing.T) {
	t.Fatal("I always fail")
}
