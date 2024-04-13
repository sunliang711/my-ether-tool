/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package script

import (
	cmd "met/cmd"
	utils "met/utils"
	"os"

	"github.com/spf13/cobra"
)

// TxCmd represents the tx command
var ScriptCmd = &cobra.Command{
	Use:   "script",
	Short: "gen script",
	Long:  `generate setup script`,
	Run:   scriptGen,
}

func init() {
	cmd.RootCmd.AddCommand(ScriptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// txCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// txCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func scriptGen(cmd *cobra.Command, args []string) {
	logger := utils.GetLogger("scriptGen")

	logger.Info().Msgf("generate add network script")

	addNetworkScriptStr := `
#!/bin/bash
if [ -z "$ID" ]; then
    echo "Missing env var ID" 1>&2
    exit 1
fi

set -e

met network add --name eth --rpc "https://mainnet.infura.io/v3/${ID}" --explorer https://etherscan.io --symbol ETH
met network add --name  polygon --rpc https://polygon-rpc.com --explorer https://polygonscan.com --symbol MATIC
met network add --name  bsc --rpc https://bsc-dataseed1.binance.org --explorer https://bscscan.com --symbol BNB
met network add --name  op --rpc https://mainnet.optimism.io --explorer https://optimistic.etherscan.io --symbol ETH
met network add --name  arbitrum --rpc https://arb1.arbitrum.io/rpc --explorer https://arbiscan.io --symbol ETH
met network add --name  arbi --rpc https://arb1.arbitrum.io/rpc --explorer https://arbiscan.io --symbol ETH
met network add --name  goerli --rpc "https://goerli.infura.io/v3/${ID}" --explorer https://goerli.etherscan.io --symbol GETH
met network add --name  sepolia --rpc "https://sepolia.infura.io/v3/${ID}" --explorer https://sepolia.etherscan.io --symbol ETH
met network add --name  ftm --rpc https://1rpc.io/ftm --explorer https://ftmscan.com --symbol FTM
met network add --name  ftmTest --rpc https://rpc.ankr.com/fantom_testnet --explorer https://testnet.ftmscan.com --symbol FTM
met network add --name  avax --rpc https://1rpc.io/avax/c --explorer https://snowtrace.io --symbol AVAX
met network add --name  bscTest --rpc https://data-seed-prebsc-1-s2.binance.org:8545 --explorer https://testnet.bscscan.com --symbol tBNB
	`

	scriptName := "addNet.sh"
	err := os.WriteFile(scriptName, []byte(addNetworkScriptStr), 0755)
	if err != nil {
		logger.Error().Err(err).Msgf("generate script: %v", scriptName)
	} else {
		logger.Info().Msgf("%v generated", scriptName)
	}
}
