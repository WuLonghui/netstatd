package main

import (
	"log"

	. "netstatd"
	"netstatd/api/server"
)

func main() {
	netstatd := NewNetstatd()
	err := netstatd.Run()
	if err != nil {
		log.Fatal(err)
	}

	server.Run(netstatd)
}
