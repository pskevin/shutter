package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// receiveCmd represents the receive command
var receiveCmd = &cobra.Command{
	Use:   "Receive",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		fmt.Println(args)
		if len(args) < 1 || len(args) > 2 {
			return errors.New("Requires at least [receiverID] and optionally [senderID] positional arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		command, receiverID := "Receive", args[0]
		req := fmt.Sprintf("%s/%s?receiverID=%s", ClusterURL, command, receiverID)
		if len(args) == 2 {
			senderID := args[1]
			req = fmt.Sprintf("%s&senderID=%s", req, senderID)
		}

		fmt.Printf("\n%s called\n", command)
		res, err := http.Get(req)
		if err != nil {
			log.Fatalln(err)
		}
		res.Body.Close()
	},
}

func init() {
	rootCmd.AddCommand(receiveCmd)
}
