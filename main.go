package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var DisplayName string = "Unset"
var ShortName string = "unset"
var Version string = "?.?.?"
var Commit string = "???????"

func main() {
	fmt.Println(DisplayName + " version v" + Version + ", build " + Commit)

	err := godotenv.Load()
	if err == nil {
		fmt.Println("Environment variables from .env loaded")
	}

	subnetString := os.Getenv("WGG_SUBNET")
	if len(subnetString) <= 0 {
		log.Fatalln("the WGG_SUBNET env var is not set or empty")
	}

	_, subnet, err := net.ParseCIDR(subnetString)
	if err != nil {
		log.Fatalln(
			"error while parsing WGG_SUBNET env var as CIDR: value '" +
				subnetString + "': " +
				err.Error(),
		)
	}

	outDir, err := InitOutDir()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("Output dir: " + outDir)
	err = CleanUpOutDir(outDir)
	if err != nil {
		log.Fatalln(err.Error())
	}

	privateNodeKey, publicNodeKey, privateClientKey, publicClientKey, err := InitKeys(outDir)
	if err != nil {
		log.Fatalln(err.Error())
	}

	nodeList, err := InitNodeList(
		privateNodeKey,
		publicNodeKey,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	PrintNodes(subnet, nodeList)

	clientList, err := InitClientList(
		privateClientKey,
		publicClientKey,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	PrintClients(subnet, clientList)

	err = GenerateNodeConfigs(
		subnet,
		outDir,
		nodeList,
		clientList,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = GenerateClientConfigs(
		subnet,
		outDir,
		nodeList,
		clientList,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("Everything is ready in " + outDir)
}

func CleanUpOutDir(outDir string) error {
	files, err := os.ReadDir(outDir)
	if err != nil {
		return errors.New("Error reading outDir: " + err.Error())
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "node.") ||
			strings.HasPrefix(file.Name(), "client.") {
			err = os.Remove(outDir + "/" + file.Name())
			if err != nil {
				return errors.New("Error removing '" + outDir + "/" + file.Name() + "': " + err.Error())
			}
		}
	}

	return nil
}

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

func InitOutDir() (string, error) {
	outDir := os.Getenv("WGG_OUT_DIR")
	if len(outDir) <= 0 {
		return "", errors.New("the WGG_OUT_DIR env var is not set or empty")
	} else if !strings.HasPrefix(outDir, "/") {
		outDir = FatalCwd() + "/" + outDir
	}

	_, err := os.Stat(outDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(outDir, 0755)
		if err != nil {
			return "", errors.New("Error creating outDir at '" + outDir + "': " + err.Error())
		}
	}

	return outDir, nil
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

func InitKeys(outDir string) (string, string, string, string, error) {
	if !IsCommandAvailable("wg") {
		return "", "", "", "", errors.New("the wg command is not available on the system, please install it")
	}

	err := InitWireGuardKeys(
		outDir+"/node.private.key",
		outDir+"/node.public.key",
	)
	if err != nil {
		return "", "", "", "", errors.New("error while creating node keys: " + err.Error())
	}

	privateNodeKey, publicNodeKey, err := LoadWireGuardKey(
		outDir+"/node.private.key",
		outDir+"/node.public.key",
	)
	if err != nil {
		return "", "", "", "", errors.New("error while loading node keys: " + err.Error())
	}

	err = InitWireGuardKeys(
		outDir+"/client.private.key",
		outDir+"/client.public.key",
	)
	if err != nil {
		return "", "", "", "", errors.New("error while creating client keys: " + err.Error())
	}

	privateClientKey, publicClientKey, err := LoadWireGuardKey(
		outDir+"/client.private.key",
		outDir+"/client.public.key",
	)
	if err != nil {
		return "", "", "", "", errors.New("error while loading client keys: " + err.Error())
	}

	return privateNodeKey, publicNodeKey, privateClientKey, publicClientKey, nil
}

func Cwd() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	return cwd, nil
}

func FatalCwd() string {
	cwd, err := Cwd()
	if err != nil {
		log.Fatalln(err)
	}
	return cwd
}

// IsCommandAvailable returns true if the command is available in the system's PATH, false otherwise.
func IsCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// GenerateWireGuardKeys generates a new WireGuard key pair.
//
// This function runs the "wg genkey" command to create a private key
// and then pipes it into the "wg pubkey" command to derive the corresponding
// public key. It returns the private key, public key, and an error if any
// command execution fails.
func GenerateWireGuardKeys() (string, string, error) {
	cmdGenKey := exec.Command("wg", "genkey")

	privateKeyBuf := &bytes.Buffer{}
	cmdGenKey.Stdout = privateKeyBuf

	if err := cmdGenKey.Run(); err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %w", err)
	}

	privateKeyString := privateKeyBuf.String()

	cmdPubKey := exec.Command("wg", "pubkey")
	cmdPubKey.Stdin = bytes.NewReader([]byte(privateKeyString))

	publicKeyBuf := &bytes.Buffer{}
	cmdPubKey.Stdout = publicKeyBuf

	if err := cmdPubKey.Run(); err != nil {
		return "", "", fmt.Errorf("failed to generate public key: %w", err)
	}

	privateKeyString = strings.TrimSpace(privateKeyString)
	publicKeyString := strings.TrimSpace(publicKeyBuf.String())

	return privateKeyString, publicKeyString, nil
}

// InitWireGuardKeys initializes a new WireGuard key pair and writes it to the
// given paths.
//
// If the private key and public key files already exist, the function does
// nothing and returns nil. Otherwise, it generates a new key pair and writes
// the private key and public key to the respective files, making sure that only
// the owner can read them.
//
// If any command execution or file I/O fails, the function returns an error.
func InitWireGuardKeys(privateKeyPath, publicKeyPath string) error {
	var keyIsMissing bool = false
	var err error
	_, err = os.Stat(privateKeyPath)
	if os.IsNotExist(err) {
		keyIsMissing = true
	}

	if !keyIsMissing {
		_, err = os.Stat(publicKeyPath)
		if os.IsNotExist(err) {
			keyIsMissing = true
		}
	}

	if !keyIsMissing {
		return nil
	}

	privateKey, publicKey, err := GenerateWireGuardKeys()
	if err != nil {
		return fmt.Errorf("error generating keys: %w", err)
	}

	err = os.WriteFile(privateKeyPath, []byte(privateKey), 0600)
	if err != nil {
		return fmt.Errorf("error writing private key to file: %w", err)
	}

	err = os.WriteFile(publicKeyPath, []byte(publicKey), 0600)
	if err != nil {
		return fmt.Errorf("error writing public key to file: %w", err)
	}

	return nil
}

func LoadWireGuardKey(privateKeyPath, publicKeyPath string) (string, string, error) {
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", "", fmt.Errorf("error reading private key: %w", err)
	} else if len(privateKey) == 0 {
		return "", "", errors.New("private key is empty")
	}

	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", "", fmt.Errorf("error reading public key: %w", err)
	} else if len(publicKey) == 0 {
		return "", "", errors.New("public key is empty")
	}

	return string(privateKey), string(publicKey), nil
}

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
	return IncrementIP(subnet.IP, node.ID+1)
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
	return IncrementIP(
		BroadcastAddress(subnet),
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

func BroadcastAddress(subnet *net.IPNet) net.IP {
	n := len(subnet.IP)
	out := make(net.IP, n)
	var m byte
	for i := 0; i < n; i++ {
		m = subnet.Mask[i] ^ 0xff
		out[i] = subnet.IP[i] | m
	}
	return out
}

func NextSubnet(subnet *net.IPNet) *net.IPNet {
	n := len(subnet.IP)
	out := BroadcastAddress(subnet)
	var c byte = 1
	for i := n - 1; i >= 0; i-- {
		out[i] = out[i] + c
		if out[i] == 0 && c > 0 {
			c = 1
		} else {
			c = 0
		}

	}
	if c == 1 {
		return nil
	}
	return &net.IPNet{IP: out.Mask(subnet.Mask), Mask: subnet.Mask}
}

func IncrementIP(ip net.IP, increment int) net.IP {
	if ip == nil {
		return nil
	}

	if increment == 0 {
		return ip.To16()
	}

	ipv4 := ip.To4()
	if ipv4 != nil {
		ip = ipv4
	}

	ipInt := new(big.Int)
	ipInt.SetBytes(ip)

	incrementInt := big.NewInt(int64(increment))

	result := new(big.Int).Add(ipInt, incrementInt)

	resultBytes := result.Bytes()

	if ip.To4() != nil {
		if len(resultBytes) > 4 {
			return nil
		}
		resultIP := make(net.IP, 4)
		copy(resultIP[4-len(resultBytes):], resultBytes)
		return resultIP
	} else {
		if len(resultBytes) > 16 {
			return nil
		}
		resultIP := make(net.IP, 16)
		copy(resultIP[16-len(resultBytes):], resultBytes)
		return resultIP
	}
}
