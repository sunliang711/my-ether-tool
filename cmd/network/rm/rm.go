package rm

import (
	"met/cmd/network"
	database "met/database"
	utils "met/utils"

	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:     "rm",
	Aliases: []string{"remove", "delete", "del"},
	Short:   "remove network",
	Long:    "remove netowrk",
	Run:     removeNetwork,
}

var (
	name *string
)

func init() {
	network.NetworkCmd.AddCommand(rmCmd)

	name = rmCmd.Flags().String("name", "", "network name")

}

func removeNetwork(cmd *cobra.Command, args []string) {
	logger := utils.GetLogger("removeNetwork")
	utils.ExitWhen(logger, *name == "", "need name\n")

	err := database.RemoveNetwork(*name)
	utils.ExitWhenErr(logger, err, "remove network error: %s", err)
}
