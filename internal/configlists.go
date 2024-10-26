package wgg

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func PrintNodes(
	subnet *net.IPNet,
	nodeList []WggNode,
) {
	fmt.Println("Nodes:")
	for _, node := range nodeList {
		fmt.Println(
			"- #" + strconv.Itoa(node.ID) +
				"| " + (*node.PubIp).String() + ":" +
				strconv.Itoa(node.Port) +
				" > " + node.WireGuardSubnetIP(subnet).String(),
		)
	}
}

func PrintClients(
	subnet *net.IPNet,
	clientList []WggClient,
) {
	fmt.Println("Clients:")
	for _, client := range clientList {
		fmt.Println(
			"- #" + strconv.Itoa(client.ID) +
				"| <client>:" +
				strconv.Itoa(client.Port) +
				" > " + client.WireGuardSubnetIP(subnet).String(),
		)
	}
}

func GenerateNodeConfigs(
	subnet *net.IPNet,
	outDir string,
	nodeList []WggNode,
	clientList []WggClient,
) error {
	var selfConf string
	var otherConfs []string
	for _, node := range nodeList {
		otherConfs = []string{}

		for _, node2 := range nodeList {
			if node.ID == node2.ID {
				selfConf = node2.WgConf(subnet, true)
			} else {
				otherConfs = append(otherConfs, node2.WgConf(subnet, false))
			}
		}

		for _, client := range clientList {
			otherConfs = append(otherConfs, client.WgConf(subnet, false))
		}

		outFile := outDir + "/node." + strconv.Itoa(node.ID) + ".wg.conf"

		err := os.WriteFile(outFile, []byte(selfConf+"\n"+strings.Join(otherConfs, "\n")), 0640)
		if err != nil {
			return errors.New("Error writing to '" + outFile + "': " + err.Error())
		}
	}

	return nil
}

func GenerateClientConfigs(
	subnet *net.IPNet,
	outDir string,
	nodeList []WggNode,
	clientList []WggClient,
) error {
	var selfConf string
	var otherConfs []string
	for _, client := range clientList {
		selfConf = client.WgConf(subnet, true)
		otherConfs = []string{}

		for _, node := range nodeList {
			otherConfs = append(otherConfs, node.WgConf(subnet, false))
		}

		outFile := outDir + "/client." + strconv.Itoa(client.ID) + ".wg.conf"
		err := os.WriteFile(outFile, []byte(selfConf+"\n"+strings.Join(otherConfs, "\n")), 0640)
		if err != nil {
			return errors.New("Error writing to '" + outFile + "': " + err.Error())
		}
	}

	return nil
}

func InitNodeList(
	privateNodeKey string,
	publicNodeKey string,
) ([]WggNode, error) {
	nodeRawDataList := []string{}

	var i int = 0
	for {
		nodeRawData := os.Getenv("WGG_NODE" + strconv.Itoa(i+1))
		if len(nodeRawData) <= 0 {
			break
		}
		nodeRawDataList = append(nodeRawDataList, nodeRawData)
		i++
	}

	nodeList := []WggNode{}

	for i, nodeRawData := range nodeRawDataList {
		node, err := NewWggNode(
			i,
			nodeRawData,
			privateNodeKey,
			publicNodeKey,
		)

		if err != nil {
			return nil, errors.New("error while creating node: " + err.Error())
		}

		nodeList = append(nodeList, node)
	}

	return nodeList, nil
}

func InitClientList(
	privateClientKey string,
	publicClientKey string,
) ([]WggClient, error) {
	clientCountString := os.Getenv("WGG_CLIENT_COUNT")
	if len(clientCountString) <= 0 {
		return nil, errors.New("the WGG_CLIENT_COUNT env var is not set or empty")
	}

	clientCount, err := strconv.Atoi(clientCountString)
	if err != nil {
		return nil, errors.New(
			"error while parsing WGG_CLIENT_COUNT as int: value '" +
				clientCountString + "': " +
				err.Error(),
		)
	} else if clientCount < 0 {
		return nil, errors.New("the WGG_CLIENT_COUNT env var must be greater than 0")
	}

	clientList := []WggClient{}
	if clientCount > 0 {
		for i := 0; i < clientCount; i++ {
			// the client priv and pub keys are already generated and loaded
			// create a new client and add it to the list

			client := NewWggClient(
				i,
				privateClientKey,
				publicClientKey,
			)

			clientList = append(clientList, client)
		}
	}

	return clientList, nil
}
