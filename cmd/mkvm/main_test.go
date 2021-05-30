package main

import "testing"

func TestGetName(t *testing.T) {
	name, err := getName()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(name)
}

func TestGetDistros(t *testing.T) {
	_, err := getDistros()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRandomMac(t *testing.T) {
	mac, err := randomMac()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(mac)
}
