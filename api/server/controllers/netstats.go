package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
	. "netstatd/namespace"
	. "netstatd/namespace/discovery"
)

type NetStatsController struct {
	beego.Controller
}

func (c *NetStatsController) ShowAll() {
	SetSuccessOutput(c.Ctx, 200, D.GetAllNetStats())
}

func (c *NetStatsController) Show() {
	var (
		namespace *Namespace
		err       error
	)

	dockerContainerId := c.GetString("docker_container_id")
	if dockerContainerId != "" {
		dockerDiscovery := NewDockerDiscovery()
		namespace, err = dockerDiscovery.GetNamespace(dockerContainerId)
		if err != nil {
			SetErrorOutput(c.Ctx, 500, err)
			return
		}
	}

	if namespace == nil {
		namespace = NewNamespace(CURRENT_NAMESPACE_PID, "host")
	}

	netStats := D.GetNetStats(namespace)

	iface := c.Ctx.Input.Param(":interface")
	netStat, ok := netStats[iface]
	if !ok {
		SetErrorOutput(c.Ctx, 404, fmt.Errorf("Interface not found"))
		return
	}
	SetSuccessOutput(c.Ctx, 200, netStat)
}
