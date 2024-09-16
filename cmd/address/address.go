package address

import (
	cmd "met/cmd"

	"github.com/spf13/cobra"
)

// AddressCmd represents the address command
var AddressCmd = &cobra.Command{
	Use:   "address",
	Short: "ethereum address related",
	Run:   nil,
}

func init() {
	cmd.RootCmd.AddCommand(AddressCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// txCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// txCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
