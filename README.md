<div align="center">

# 🕷️🕸️ wgg

![CI/CD](https://github.com/CoreUnit-NET/wgg/actions/workflows/go-bin-release.yml/badge.svg)
![CI/CD](https://github.com/CoreUnit-NET/wgg/actions/workflows/go-test-build.yml/badge.svg)  
![MIT](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/dynamic/json?color=green&label=watchers&query=watchers&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2FCoreUnit-NET%2Fwgg)
![](https://img.shields.io/badge/dynamic/json?color=yellow&label=stars&query=stargazers_count&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2FCoreUnit-NET%2Fwgg)
![](https://img.shields.io/badge/dynamic/json?color=navy&label=forks&query=forks&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2FCoreUnit-NET%2Fwgg)

</div>

`wgg` simply generates server p2p node and client configurations.

<details><summary><strong>Configuration</strong></summary>

## Configuration

WGG is fully configured via env vars or a .env file:

```bash
WGG_SUBNET=10.10.10.0/24
WGG_NODE1=<node1-ip>:55333 # if you subsequently adjust or add a node, all configs must be adjusted
WGG_NODE2=<node2-ip>:55333
WGG_NODE3=<node3-ip>:55333
WGG_CLIENT_COUNT=10 #tip: choose a number that is sufficient for users in the long term, whereby all node configs must be updated for each new user
WGG_OUT_DIR=config
```

</details>

<details><summary><strong>User Guide</strong></summary>

# User Guide

## Requirements

Linux- or macos-like systems with `go` or `wget & tar` installed.

## Getting Started

Start the latest repo version directly without leaving stuff in the current working dir:

```sh
go run github.com/CoreUnit-NET/wgg@latest
```

## Quick help

```sh
go run github.com/CoreUnit-NET/wgg@latest -h
```

## Install via go

###### _For this section go is required, check out the [install go guide](#install-go)._

```sh
go install github.com/CoreUnit-NET/wgg@latest
```

## Install via wget

```sh
export CUSTOM_BIN_DIR="/usr/local/bin" # <- change if needed
export CUSTOM_VERSION="" # <- set latest version here

rm -rf $CUSTOM_BIN_DIR/wgg
wget https://github.com/CoreUnit-NET/wgg/releases/download/v$CUSTOM_VERSION/wgg-v$CUSTOM_VERSION-linux-amd64.tar.gz -O /tmp/wgg.tar.gz
tar -xzvf /tmp/wgg.tar.gz -C $CUSTOM_BIN_DIR/ wgg
rm /tmp/wgg.tar.gz
```

# Build

## Build requirements

To build, you need to install go.
The required go version is in the `go.mod` file.

## Build Instructions

###### _For this section go is required, check out the [install go guide](#install-go)._

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

</details>

<details><summary><strong>Development</strong></summary>

# Development

###### _For this section go is required, check out the [install go guide](#install-go)._

This part is work in progress, I want to use 'AIR' as auto-reload tool:

```sh
make dev #WIP
```

## Install go

The required go version for this project is in the `go.mod` file.

To install and update go, I can recommend the following repo:

```sh
git clone git@github.com:udhos/update-golang.git golang-updater
cd golang-updater
sudo ./update-golang.sh
```

</details>

<div align="center">

# 🤝 Contributing

Contributions to this project are welcome!  
Follow the [CONTRIBUTING.md](CONTRIBUTING.md) for more infos.

# ⚠️ Disclaimer

This project is provided without warranties.

# 📜 License

Licensed under the [MIT license](LICENSE).

<a href="https://discord.coreunit.net">
    <img alt="CoreUnit.NET Discord Banner" src="https://discord.com/api/guilds/422136748294930443/widget.png?style=banner2">
</a>

</div>