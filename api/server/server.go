package server

import (
	. "netstatd"
	"netstatd/api/server/controllers"
	_ "netstatd/api/server/routers"

	"github.com/astaxie/beego"
)

func Run(netstatd *Netstatd) {
	controllers.Init(netstatd)
	beego.Run()
}
