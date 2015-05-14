package routers

import (
	"netstatd/api/server/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSRouter("/netstats", &controllers.NetStatsController{}, "get:ShowAll"),
	)
	beego.AddNamespace(ns)
}
