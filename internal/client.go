package wgg

import (
	"fmt"
	"net"

	"coreunit.net/wgg/lib/netutils"
)

type WggClient struct {
	ID         int
	Port       int
	PrivateKey string
	PublicKey  string
}

func NewWggClient(
	id int,
	privateKey string,
	publicKey string,
) WggClient {
	return WggClient{
		ID:         id,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}

func (client WggClient) WireGuardSubnetIP(subnet *net.IPNet) net.IP {
	return netutils.IncrementIP(
		netutils.BroadcastAddress(subnet),
		-(client.ID + 1),
	)
}

func (client WggClient) WgConf(subnet *net.IPNet, forSelf bool) string {
	if forSelf {
		return fmt.Sprintf(
			"[Interface]\n"+
				"PrivateKey = %s\n"+
				"Address = %s\n"+
				"",
			client.PrivateKey,
			client.WireGuardSubnetIP(subnet),
		)
	} else {
		return fmt.Sprintf(
			"[Peer]\n"+
				"PublicKey = %s\n"+
				"AllowedIPs = %s\n"+
				"",
			client.PublicKey,
			client.WireGuardSubnetIP(subnet),
		)
	}
}
