package wgg

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"coreunit.net/wgg/lib/netutils"
)

type WggNode struct {
	ID         int
	PubIp      *net.IP
	Port       int
	PrivateKey string
	PublicKey  string
}

// WireGuardSubnetIP returns an IP address in the given subnet that is
// appropriate for the current node to use as its WireGuard IP address.
//
// The returned IP address is the given subnet's IP address incremented by
// the node's ID plus one.
func (node WggNode) WireGuardSubnetIP(subnet *net.IPNet) net.IP {
	return netutils.IncrementIP(subnet.IP, node.ID+1)
}

// WgConf generates a WireGuard config string for the current node.
//
// If forSelf is true, the generated config is for the current node itself.
// Otherwise, it is for a peer of the current node.
//
// The generated config does not include the [Interface] section if forSelf is
// false.
func (node WggNode) WgConf(subnet *net.IPNet, forSelf bool) string {
	if forSelf {
		return fmt.Sprintf(
			"[Interface]\n"+
				"Address = %s\n"+
				"PrivateKey = %s\n"+
				"ListenPort = %d\n"+
				"",
			node.WireGuardSubnetIP(subnet),
			node.PrivateKey,
			node.Port,
		)
	} else {
		return fmt.Sprintf(
			"[Peer]\n"+
				"PublicKey = %s\n"+
				"AllowedIPs = %s\n"+
				"Endpoint = %s:%d\n"+
				"",
			node.PublicKey,
			node.WireGuardSubnetIP(subnet),
			*node.PubIp,
			node.Port,
		)
	}
}

func NewWggNode(
	id int,
	rawData string,
	privateKey string,
	publicKey string,
) (WggNode, error) {
	host, portStr, err := net.SplitHostPort(rawData)
	if err != nil {
		return WggNode{}, errors.New(
			"general invalid raw node data: '" +
				rawData + "': " +
				err.Error(),
		)
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return WggNode{}, errors.New(
			"invalid ip in raw node data: '" +
				rawData + "': " +
				err.Error(),
		)
	}

	// Parse the port
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return WggNode{}, errors.New(
			"invalid port in raw node data: '" +
				rawData + "': " +
				err.Error(),
		)
	}

	return WggNode{
		ID:         id,
		PubIp:      &ip,
		Port:       port,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}
