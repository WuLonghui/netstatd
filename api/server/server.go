package server

import (
	"log"
	. "netstatd"
	"netstatd/api/server/controllers"
	_ "netstatd/api/server/routers"

	"github.com/astaxie/beego"
	"github.com/pivotal-golang/localip"
)

func Run(netstatd *Netstatd) {
	controllers.Init(netstatd)

	//important: server must listen on localIp, or it will block when capturing package
	localIp, err := localip.LocalIP()
	if err != nil {
		log.Fatal(err)
	}
	beego.HttpAddr = localIp

	beego.Run()
}
