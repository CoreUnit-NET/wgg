package wgg

import (
	"fmt"
	"net"
)

func GenWgClientConfPart(
	target WggTarget,
	keyDir string,
	subnet *net.IPNet,
	forTargetID string,
) (string, error) {
	ones, _ := subnet.Mask.Size()

	if target.TargetID() == forTargetID {
		privateKey, _, err := InitWireGuardKeyPair(
			keyDir+"/"+target.TargetID()+".key",
			keyDir+"/"+target.TargetID()+".pub",
		)
		if err != nil {
			return "", err
		}

		if target.IsNode() {
			return fmt.Sprintf(
				"[Interface]\n"+
					"Address = %s/%d\n"+
					"PrivateKey = %s\n"+
					"ListenPort = %d\n",
				target.WireGuardSubnetIP(subnet),
				ones,
				privateKey,
				target.NodePort(),
			), nil
		} else {
			return fmt.Sprintf(
				"[Interface]\n"+
					"PrivateKey = %s\n"+
					"Address = %s/%d\n",
				privateKey,
				target.WireGuardSubnetIP(subnet),
				ones,
			), nil
		}
	} else {
		_, publicKey, err := InitWireGuardKeyPair(
			keyDir+"/"+target.TargetID()+".key",
			keyDir+"/"+target.TargetID()+".pub",
		)
		if err != nil {
			return "", err
		}

		if target.IsNode() {
			return fmt.Sprintf(
				"[Peer]\n"+
					"PublicKey = %s\n"+
					"AllowedIPs = %s/32\n"+
					"Endpoint = %s:%d\n",
				publicKey,
				target.WireGuardSubnetIP(subnet),
				target.NodePubIp(),
				target.NodePort(),
			), nil
		} else {
			return fmt.Sprintf(
				"[Peer]\n"+
					"PublicKey = %s\n"+
					"AllowedIPs = %s/32\n",
				publicKey,
				target.WireGuardSubnetIP(subnet),
			), nil
		}
	}
}
