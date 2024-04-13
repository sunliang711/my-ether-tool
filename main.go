/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	cmd "met/cmd"
	setup "met/setup"

	_ "met/cmd/account"
	_ "met/cmd/account/add"
	_ "met/cmd/account/current"
	_ "met/cmd/account/list"
	_ "met/cmd/account/new"
	_ "met/cmd/account/rm"
	_ "met/cmd/account/switch"

	_ "met/cmd/contract"
	_ "met/cmd/contract/read"
	_ "met/cmd/contract/write"

	_ "met/cmd/erc20"
	_ "met/cmd/erc20/allowance"
	_ "met/cmd/erc20/approve"
	_ "met/cmd/erc20/balanceOf"
	_ "met/cmd/erc20/decimals"
	_ "met/cmd/erc20/name"
	_ "met/cmd/erc20/symbol"
	_ "met/cmd/erc20/totalSupply"
	_ "met/cmd/erc20/transfer"
	_ "met/cmd/erc20/transferFrom"

	_ "met/cmd/codec"
	_ "met/cmd/codec/decode"
	_ "met/cmd/codec/encode"

	_ "met/cmd/network"
	_ "met/cmd/network/add"
	_ "met/cmd/network/current"
	_ "met/cmd/network/list"
	_ "met/cmd/network/rm"
	_ "met/cmd/network/switch"

	_ "met/cmd/tx"
	_ "met/cmd/tx/offsign"
	_ "met/cmd/tx/send"

	_ "met/cmd/script"
)

func main() {
	setup.SetupDb()

	cmd.Execute()
}
