package main

import (
	"fmt"
	"patdown/common"
)

func main() {
	common.LoadArgs()
	var servers []string

	common.Banner()

	autodetect := common.Params.Domain != ""
	if autodetect {
		if servers = common.PullNS(common.Params.Domain); len(servers) == 0 {
			common.Fatal("no nameservers found for " + common.Params.Domain)
		}
		common.Info(fmt.Sprintf("retrieved %s%d%s nameservers for %s", common.ColorGreen, len(servers), common.ColorReset, common.Params.Domain))
	} else if len(common.Params.Nservers) > 0 {
		servers = common.Params.Nservers
	} else {
		common.Fatal("provide a domain or nameservers to target")
	}

	if !common.NeutralReq() {
		common.Fatal("neutral dns check failed, are you on a dirty box or vpn?")
	}

	valid := common.ParseNS(servers)
	if len(valid) == 0 {
		common.Fatal("no servers responded to trial probes, they're either down or they don't like your IP")
	}

	common.Takeoff(valid)
}
