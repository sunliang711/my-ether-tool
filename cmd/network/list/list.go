package list

import (
	"met/cmd/network"
	database "met/database"
	utils "met/utils"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"show"},
	Short:   "list network",
	Long:    "list network",
	Run:     listNetwork,
}

var name *string

func init() {
	network.NetworkCmd.AddCommand(listCmd)

	name = listCmd.Flags().String("name", "", "show specify network instead of all networks")
}

func listNetwork(cmd *cobra.Command, args []string) {
	if *name != "" {
		net, err := database.QueryNetwork(*name)
		utils.ExitWhenError(err, "query network: %s error: %s", *name, err)
		network.ShowNetwork(net)
	} else {

		networks, err := database.QueryAllNetworks()
		utils.ExitWhenError(err, "query networks error: %s", err)

		for i := range networks {
			network.ShowNetwork(networks[i])
		}
	}

}
