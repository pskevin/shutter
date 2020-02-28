package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// createNodeCmd represents the createNode command
var createNodeCmd = &cobra.Command{
	Use:   "CreateNode",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("Requires exactly [nodeID] [initAmount] positional arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		command, nodeID, initAmount := "CreateNode", args[0], args[1]
		fmt.Printf("\n%s called\n", command)
		res, err := http.Get(fmt.Sprintf("%s/%s?nodeID=%s&initAmount=%s", ClusterURL, command, nodeID, initAmount))
		if err != nil {
			log.Fatalln(err)
		}
		res.Body.Close()
	},
}

func init() {
	rootCmd.AddCommand(createNodeCmd)
}
