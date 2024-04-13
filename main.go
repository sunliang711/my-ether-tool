/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"my-ether-tool/cmd"
	"my-ether-tool/setup"

	_ "my-ether-tool/cmd/account"
	_ "my-ether-tool/cmd/account/add"
	_ "my-ether-tool/cmd/account/current"
	_ "my-ether-tool/cmd/account/list"
	_ "my-ether-tool/cmd/account/new"
	_ "my-ether-tool/cmd/account/rm"
	_ "my-ether-tool/cmd/account/switch"

	_ "my-ether-tool/cmd/contract"
	_ "my-ether-tool/cmd/contract/read"
	_ "my-ether-tool/cmd/contract/write"

	_ "my-ether-tool/cmd/erc20"
	_ "my-ether-tool/cmd/erc20/allowance"
	_ "my-ether-tool/cmd/erc20/approve"
	_ "my-ether-tool/cmd/erc20/balanceOf"
	_ "my-ether-tool/cmd/erc20/decimals"
	_ "my-ether-tool/cmd/erc20/name"
	_ "my-ether-tool/cmd/erc20/symbol"
	_ "my-ether-tool/cmd/erc20/totalSupply"
	_ "my-ether-tool/cmd/erc20/transfer"
	_ "my-ether-tool/cmd/erc20/transferFrom"

	_ "my-ether-tool/cmd/codec"
	_ "my-ether-tool/cmd/codec/decode"
	_ "my-ether-tool/cmd/codec/encode"

	_ "my-ether-tool/cmd/network"
	_ "my-ether-tool/cmd/network/add"
	_ "my-ether-tool/cmd/network/current"
	_ "my-ether-tool/cmd/network/list"
	_ "my-ether-tool/cmd/network/rm"
	_ "my-ether-tool/cmd/network/switch"

	_ "my-ether-tool/cmd/tx"
	_ "my-ether-tool/cmd/tx/offsign"
	_ "my-ether-tool/cmd/tx/send"
)

func main() {
	setup.SetupDb()

	cmd.Execute()
}
