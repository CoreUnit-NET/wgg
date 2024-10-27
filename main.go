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

	// err = Test()
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }

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

	outDir, keyDir, err := wgg.InitOutDir()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("Output dir: " + outDir)
	err = wgg.CleanUpOutDir(outDir)
	if err != nil {
		log.Fatalln(err.Error())
	}

	nodeList, err := wgg.InitNodeList()
	if err != nil {
		log.Fatalln(err.Error())
	}

	wgg.PrintNodes(subnet, nodeList)

	clientList, err := wgg.InitClientList()
	if err != nil {
		log.Fatalln(err.Error())
	}

	wgg.PrintClients(subnet, clientList)

	err = wgg.GenerateNodeConfigs(
		subnet,
		outDir,
		keyDir,
		nodeList,
		clientList,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = wgg.GenerateClientConfigs(
		subnet,
		outDir,
		keyDir,
		nodeList,
		clientList,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("Everything is ready in " + outDir)
}

// func Test() error {
// 	homeDir := os.Getenv("HOME")
// 	if len(homeDir) == 0 {
// 		return errors.New("HOME env var is not set")
// 	}

// 	privateKey, err := os.ReadFile(homeDir + "/.ssh/id_rsa")
// 	if err != nil {
// 		return errors.New("error reading private key: " + err.Error())
// 	}

// 	config := &sftputils.SshConfig{
// 		Host:       "37.131.245.163",
// 		User:       "cunet",
// 		PrivateKey: string(privateKey),
// 	}

// 	return sftputils.HandleSftp(
// 		config,
// 		func(sftp *sftp.Client, session *ssh.Session) error {
// 			path, err := sftputils.JoinPath(sftp, "~/.ssh/known_hosts")
// 			if err != nil {
// 				return errors.New("error joining path: " + err.Error())
// 			}
// 			fmt.Println(path)

// 			f, err := sftp.Create("hello.txt")
// 			if err != nil {
// 				return errors.New("error creating file: " + err.Error())
// 			}
// 			if _, err := f.Write([]byte("Hello world!")); err != nil {
// 				return errors.New("error writing to file: " + err.Error())
// 			}
// 			f.Close()

// 			fi, err := sftp.Lstat("hello.txt")
// 			if err != nil {
// 				return errors.New("error stating file: " + err.Error())
// 			}
// 			fmt.Println(fi)

// 			// var stdin io.WriteCloser
// 			// var stdout, stderr io.Reader

// 			// stdin, err = session.StdinPipe()
// 			// if err != nil {
// 			// 	fmt.Println(err.Error())
// 			// }

// 			// stdout, err = session.StdoutPipe()
// 			// if err != nil {
// 			// 	fmt.Println(err.Error())
// 			// }

// 			// stderr, err = session.StderrPipe()
// 			// if err != nil {
// 			// 	fmt.Println(err.Error())
// 			// }

// 			// wr := make(chan []byte, 10)

// 			// go func() {
// 			// 	for d := range wr {
// 			// 		_, err := stdin.Write(d)
// 			// 		if err != nil {
// 			// 			fmt.Println(err.Error())
// 			// 		}
// 			// 	}
// 			// }()

// 			// go func() {
// 			// 	scanner := bufio.NewScanner(stdout)
// 			// 	for {
// 			// 		if tkn := scanner.Scan(); tkn {
// 			// 			rcv := scanner.Bytes()

// 			// 			raw := make([]byte, len(rcv))
// 			// 			copy(raw, rcv)

// 			// 			fmt.Println(string(raw))
// 			// 		} else {
// 			// 			if scanner.Err() != nil {
// 			// 				fmt.Println(scanner.Err())
// 			// 			} else {
// 			// 				fmt.Println("io.EOF")
// 			// 			}
// 			// 			return
// 			// 		}
// 			// 	}
// 			// }()

// 			// go func() {
// 			// 	scanner := bufio.NewScanner(stderr)

// 			// 	for scanner.Scan() {
// 			// 		fmt.Println(scanner.Text())
// 			// 	}
// 			// }()

// 			// err = session.Shell()
// 			// if err != nil {
// 			// 	fmt.Println(err.Error())
// 			// }

// 			// for {
// 			// 	fmt.Println("$")

// 			// 	scanner := bufio.NewScanner(os.Stdin)
// 			// 	scanner.Scan()
// 			// 	text := scanner.Text()

// 			// 	wr <- []byte(text + "\n")
// 			// }

// 			return nil
// 		},
// 	)
// }
