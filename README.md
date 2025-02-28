# WireGuardGenerator

*This is a prove of concept and public test!*  
Generates basic p2p node and client configs with private and public keys.
You need to define env vars.

Also loads .env files.

Env vars:
```bash
WGG_SUBNET=10.10.10.0/24
WGG_NODE1=<node1-ip>:55333 # if you subsequently adjust or add a node, all configs must be adjusted
WGG_NODE2=<node2-ip>:55333
WGG_NODE3=<node3-ip>:55333
WGG_CLIENT_COUNT=10 #tip: choose a number that is sufficient for users in the long term, whereby all node configs must be updated for each new user
WGG_OUT_DIR=config
```

# Table of Contents
- [WireGuardGenerator](#wireguardgenerator)
- [Table of Contents](#table-of-contents)
- [Getting Started](#getting-started)
  - [Requirements](#requirements)
  - [Install via go](#install-via-go)
  - [Install via wget](#install-via-wget)
  - [Build](#build)
  - [Install go](#install-go)
- [Contributing](#contributing)
- [License](#license)
- [Disclaimer](#disclaimer)

# Getting Started

## Requirements
None windows system with `go` or `wget & tar` installed.

## Install via go
###### *For this section go is required, check out the [install go guide](#install-go).*

```sh
go install https://github.com/CoreUnit-NET/wgg
```

## Install via wget
```sh
BIN_DIR="/usr/local/bin"
WGG_VERSION="1.3.3"

rm -rf $BIN_DIR/wgg
wget https://github.com/CoreUnit-NET/wgg/releases/download/v$WGG_VERSION/wgg-v$WGG_VERSION-linux-amd64.tar.gz -O /tmp/wgg.tar.gz
tar -xzvf /tmp/wgg.tar.gz -C $BIN_DIR/ wgg
rm /tmp/wgg.tar.gz
```

## Build
###### *For this section go is required, check out the [install go guide](#install-go).*

Clone the repo:
```sh
git clone https://github.com/CoreUnit-NET/wgg.git
cd wgg
```

Build the wgg binary from source code:
```sh
make build
./wgg
```

## Install go
The required go version for this project is in the `go.mod` file.

To install and update go, I can recommend the following repo:
```sh
git clone git@github.com:udhos/update-golang.git golang-updater
cd golang-updater
sudo ./update-golang.sh
```

# Contributing
Contributions to this project are welcome!  
Interested users can refer to the guidelines provided in the [CONTRIBUTING.md](CONTRIBUTING.md) file to contribute to the project and help improve its functionality and features.

# License
This project is licensed under the [MIT license](LICENSE), providing users with flexibility and freedom to use and modify the software according to their needs.

# Disclaimer
This project is provided without warranties.  
Users are advised to review the accompanying license for more information on the terms of use and limitations of liability.
