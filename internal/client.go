package wgg

import (
	"fmt"
	"net"

	"coreunit.net/wgg/lib/netutils"
)

type WggClient struct {
	WggTarget

	ID int
}

// NewWggClient returns a new WggClient.
//
// The privateKey and publicKey args are unused at the moment, but are
// reserved for future use.
func NewWggClient(
	id int,
) WggClient {
	return WggClient{
		ID: id,
	}
}

// TargetID returns the ID of the target client as a string, prefixed with "C".
func (client WggClient) TargetID() string {
	return fmt.Sprintf("c%d", client.ID)

}

func (client WggClient) IsNode() bool {
	return false
}

func (client WggClient) NodePort() int {
	return -1
}

func (client WggClient) NodePubIp() *net.IP {
	return nil
}

// WireGuardSubnetIP returns an IP address in the given subnet that is
// appropriate for the current client to use as its WireGuard IP address.
//
// The returned IP address is the given subnet's broadcast address decremented
// by the client's ID plus one.
func (client WggClient) WireGuardSubnetIP(subnet *net.IPNet) net.IP {
	return netutils.IncrementIP(
		netutils.BroadcastAddress(subnet),
		-(client.ID + 1),
	)
}
