package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
)

type NetStatsController struct {
	beego.Controller
}

func (c *NetStatsController) Show() {
	pid, err := c.GetInt("pid")
	if err != nil {
		SetErrorOutput(c.Ctx, 400, fmt.Errorf("Invaild parameters"))
		return
	}

	d, ok := D.NS[pid]
	if !ok {
		SetErrorOutput(c.Ctx, 404, fmt.Errorf("Namespace not in stat"))
		return
	}

	iface := c.Ctx.Input.Param(":interface")
	netStat, ok := d.NetStats[iface]
	if !ok {
		SetErrorOutput(c.Ctx, 404, fmt.Errorf("Interface not found"))
		return
	}
	SetSuccessOutput(c.Ctx, 200, netStat)
}

func (c *NetStatsController) ShowAll() {
	pid, err := c.GetInt("pid")
	if err != nil {
		SetErrorOutput(c.Ctx, 400, fmt.Errorf("Invaild parameters"))
		return
	}

	d, ok := D.NS[pid]
	if !ok {
		SetErrorOutput(c.Ctx, 404, fmt.Errorf("Namespace not in stat"))
		return
	}

	SetSuccessOutput(c.Ctx, 200, d.NetStats)
}

func (c *NetStatsController) Create() {
	type CreateRequest struct {
		Pid int `json:"pid"`
	}

	var createRequest CreateRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &createRequest)
	if err != nil {
		SetErrorOutput(c.Ctx, 400, err)
		return
	}

	_, ok := D.NS[createRequest.Pid]
	if ok {
		SetErrorOutput(c.Ctx, 400, fmt.Errorf("Namespace already in stat"))
		return
	}

	err = D.AddNameSpaceStat(createRequest.Pid)
	if err != nil {
		SetErrorOutput(c.Ctx, 400, err)
		return
	}

	SetSuccessOutput(c.Ctx, 201, nil)
}
