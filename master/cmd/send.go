package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "Send",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 3 {
			return errors.New("Requires exactly [senderID] [receiverID] [amount] positional arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		command, senderID, receiverID, amount := "Send", args[0], args[1], args[2]
		fmt.Printf("\n%s called\n", command)
		res, err := http.Get(fmt.Sprintf("%s/%s?senderID=%s&receiverID=%s&amount=%s", ClusterURL, command, senderID, receiverID, amount))
		if err != nil {
			log.Fatalln(err)
		}
		res.Body.Close()
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
