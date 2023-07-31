/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"my-ether-tool/cmd"

	_ "my-ether-tool/cmd/account"
	_ "my-ether-tool/cmd/account/add"
	_ "my-ether-tool/cmd/account/create"
	_ "my-ether-tool/cmd/account/current"
	_ "my-ether-tool/cmd/account/list"
	_ "my-ether-tool/cmd/account/rm"
	_ "my-ether-tool/cmd/account/switch"

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

	"my-ether-tool/database"
)

func main() {
	database.InitDB("error", "wallet.db")

	cmd.Execute()
}
