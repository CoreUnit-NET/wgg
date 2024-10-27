package wgg

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"coreunit.net/wgg/lib/netutils"
)

type WggTarget interface {
	TargetID() string
	IsNode() bool
	WireGuardSubnetIP(*net.IPNet) net.IP
	NodePort() int
	NodePubIp() *net.IP
}

type WggNode struct {
	WggTarget

	ID    int
	PubIp *net.IP
	Port  int
}

// NewWggNode parses a raw node data string into a WggNode.
//
// The raw node data string should be in the format of "<ip>:<port>".
//
// The returned WggNode's ID is set to the given ID argument.
//
// The returned WggNode's PubIp is set to the parsed IP address.
// The returned WggNode's Port is set to the parsed port number.
//
// If the raw node data is invalid, an error is returned.
func NewWggNode(
	id int,
	rawData string,
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
				rawData + "'",
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
		ID:    id,
		PubIp: &ip,
		Port:  port,
	}, nil
}

// TargetID returns the node's target ID, which is in the format of "N<id>".
func (node WggNode) TargetID() string {
	return fmt.Sprintf("n%d", node.ID)
}

// IsNode returns true, indicating that this WggTarget is a node.
func (node WggNode) IsNode() bool {
	return true
}

// WireGuardSubnetIP returns an IP address in the given subnet that is
// appropriate for the current node to use as its WireGuard IP address.
//
// The returned IP address is the given subnet's IP address incremented by
// the node's ID plus one.
func (node WggNode) WireGuardSubnetIP(subnet *net.IPNet) net.IP {
	return netutils.IncrementIP(subnet.IP, node.ID+1)
}

// NodePort returns the port number for the current node to use as its
// WireGuard port.
func (node WggNode) NodePort() int {
	return node.Port
}

// NodePubIp returns the public IP address for the current node.
func (node WggNode) NodePubIp() *net.IP {
	return node.PubIp
}
