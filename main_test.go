package main

import (
	"regexp"
	"testing"
)

func TestDefaultPortNumber(t *testing.T) {
	address := getListeningAddress()
	match, err := regexp.Match(":\\d+", []byte(address))
	if err != nil {
		t.Error(err)
	}

	if !match {
		t.Fatalf("Bad pattern of listen address %s", address)
	}

}
