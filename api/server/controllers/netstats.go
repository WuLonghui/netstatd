package controllers

import (
	"github.com/astaxie/beego"
)

type NetStatsController struct {
	beego.Controller
}

func (c *NetStatsController) ShowAll() {
	SetSuccessOutput(c.Ctx, 200, D.NetStats)
}
