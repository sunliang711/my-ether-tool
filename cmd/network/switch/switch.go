package switchNetwork

import (
	"met/cmd/network"
	database "met/database"
	utils "met/utils"

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
	logger := utils.GetLogger("switchCurrentNetwork")

	utils.ExitWhen(logger, *name == "", "need name")

	err := database.SwitchNetwork(*name)
	utils.ExitWhenErr(logger, err, "switch network error: %s", err)
}
