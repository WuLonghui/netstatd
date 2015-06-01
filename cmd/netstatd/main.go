package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	. "netstatd"

	"netstatd/api/server"
)

var statTarget = flag.String(
	"statTarget",
	"host",
	"the target of netstatd(host,docker)",
)

func main() {

	flag.Parse()

	netstatd := NewNetstatd()

	go server.Run(netstatd)

	//waitting for server startup completely
	time.Sleep(1 * time.Second)

	err := netstatd.Run(*statTarget)
	if err != nil {
		log.Fatal(err)
	}

	killChan := make(chan os.Signal)
	signal.Notify(killChan, os.Kill, os.Interrupt)
	for {
		select {
		case <-killChan:
			log.Println("shutting down")
			return
		}
	}
}
