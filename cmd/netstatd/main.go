package main

import (
	. "netstatd"
	"netstatd/api/server"
)

func main() {
	netstatd := NewNetstatd()
	netstatd.Run()

	server.Run(netstatd)
}
