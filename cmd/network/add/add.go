package add

import (
	"my-ether-tool/cmd/network"
	"my-ether-tool/database"
	"my-ether-tool/utils"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add network",
	Long:  "add network",
	Run:   addNetwork,
}

var (
	name     *string
	rpc      *string
	symbol   *string
	explorer *string
)

func init() {
	network.NetworkCmd.AddCommand(addCmd)

	name = addCmd.Flags().String("name", "", "network name")
	rpc = addCmd.Flags().String("rpc", "", "network rpc")
	symbol = addCmd.Flags().String("symbol", "", "native token symbo,eg: ETH BNB")
	explorer = addCmd.Flags().String("explorer", "", "network explorer")
}

func addNetwork(cmd *cobra.Command, args []string) {
	var (
		err error
	)

	utils.ExitWithMsgWhen(*name == "", "need name\n")
	utils.ExitWithMsgWhen(*rpc == "", "need rpc\n")
	utils.ExitWithMsgWhen(*symbol == "", "need symbol\n")

	network := database.Network{
		Name:     *name,
		Rpc:      *rpc,
		Symbol:   *symbol,
		Explorer: *explorer,
		Current:  false,
	}
	err = database.AddNetwork(&network)
	utils.ExitWhenError(err, "Add netowrk error: %s\n", err)
}
