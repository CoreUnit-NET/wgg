# WireGuardGenerator
![CI/CD](https://github.com/noblemajo/wgg/actions/workflows/go-bin-release.yml/badge.svg)
![CI/CD](https://github.com/noblemajo/wgg/actions/workflows/go-test-build.yml/badge.svg)  
![MIT](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/dynamic/json?color=green&label=watchers&query=watchers&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2Fnoblemajo%2Fwgg)
![](https://img.shields.io/badge/dynamic/json?color=yellow&label=stars&query=stargazers_count&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2Fnoblemajo%2Fwgg)
![](https://img.shields.io/badge/dynamic/json?color=navy&label=forks&query=forks&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2Fnoblemajo%2Fwgg)


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
go install https://github.com/NobleMajo/wgg
```

## Install via wget
```sh
BIN_DIR="/usr/local/bin"
WGG_VERSION="1.3.3"

rm -rf $BIN_DIR/wgg
wget https://github.com/NobleMajo/wgg/releases/download/v$WGG_VERSION/wgg-v$WGG_VERSION-linux-amd64.tar.gz -O /tmp/wgg.tar.gz
tar -xzvf /tmp/wgg.tar.gz -C $BIN_DIR/ wgg
rm /tmp/wgg.tar.gz
```

## Build
###### *For this section go is required, check out the [install go guide](#install-go).*

Clone the repo:
```sh
git clone https://github.com/NobleMajo/wgg.git
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
