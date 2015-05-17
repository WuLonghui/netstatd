package routers

import (
	"netstatd/api/server/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSRouter("/netstats", &controllers.NetStatsController{}, "get:ShowAll"),
		beego.NSRouter("/netstats/:interface", &controllers.NetStatsController{}, "get:Show"),
		beego.NSRouter("/netstats", &controllers.NetStatsController{}, "post:Create"),
	)
	beego.AddNamespace(ns)
}
