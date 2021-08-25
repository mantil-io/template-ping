package main

import (
	"log"
	"os/exec"
)

var apiURL = ""

func init() {
	out, err := exec.Command("mantil", "env", "-u").Output()
	if err != nil {
		log.Fatalf("can't find api url, execute of `mantil env -u` faild %v", err)
	}
	apiURL = string(out)
}
