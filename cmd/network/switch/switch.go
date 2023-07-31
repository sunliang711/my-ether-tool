package switchNetwork

import (
	"my-ether-tool/cmd/network"
	"my-ether-tool/database"
	"my-ether-tool/utils"

	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "switch current network",
	Long:  "switch current network",
	Run:   switchCurrentNetwork,
}

var (
	name *string
)

func init() {
	network.NetworkCmd.AddCommand(switchCmd)

	name = switchCmd.Flags().String("name", "", "network name")

}

func switchCurrentNetwork(cmd *cobra.Command, args []string) {
	utils.ExitWithMsgWhen(*name == "", "need name\n")

	err := database.SwitchNetwork(*name)
	utils.ExitWhenError(err, "switch network error: %s\n", err)
}
