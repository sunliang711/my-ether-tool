package current

import (
	"met/cmd/network"
	database "met/database"
	utils "met/utils"

	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "show current network",
	Long:  "show current network",
	Run:   showCurrentNetwork,
}

func init() {
	network.NetworkCmd.AddCommand(currentCmd)
}

func showCurrentNetwork(cmd *cobra.Command, args []string) {
	current, err := database.CurrentNetwork()
	utils.ExitWhenError(err, "show current entwork error: %s", err)

	network.ShowNetwork(current)

}
