package main

import (
	"fmt"
	"log"
	"net"
	"os"

	wgg "coreunit.net/wgg/internal"
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

	outDir, err := wgg.InitOutDir()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("Output dir: " + outDir)
	err = wgg.CleanUpOutDir(outDir)
	if err != nil {
		log.Fatalln(err.Error())
	}

	privateNodeKey, publicNodeKey, privateClientKey, publicClientKey, err := wgg.InitAllKeys(outDir)
	if err != nil {
		log.Fatalln(err.Error())
	}

	nodeList, err := wgg.InitNodeList(
		privateNodeKey,
		publicNodeKey,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	wgg.PrintNodes(subnet, nodeList)

	clientList, err := wgg.InitClientList(
		privateClientKey,
		publicClientKey,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	wgg.PrintClients(subnet, clientList)

	err = wgg.GenerateNodeConfigs(
		subnet,
		outDir,
		nodeList,
		clientList,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = wgg.GenerateClientConfigs(
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
