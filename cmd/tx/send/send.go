package tx

import (
	"my-ether-tool/cmd/tx"

	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "send tx",
	Long:  "send transaction",
	Run:   sendTransaction,
}

func sendTransaction(cmd *cobra.Command, args []string) {

}

func init() {
	tx.TxCmd.AddCommand(sendCmd)

	// sendCmd.Flags()
}
