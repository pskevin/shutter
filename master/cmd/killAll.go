package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// killAllCmd represents the killAll command
var killAllCmd = &cobra.Command{
	Use:   "KillAll",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return errors.New("Requires exactly zero positional arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		command := "KillAll"
		fmt.Printf("\n%s called\n", command)
		res, err := http.Get(fmt.Sprintf("%s/%s", ClusterURL, command))
		if err != nil {
			log.Fatalln(err)
		}
		res.Body.Close()
	},
}

func init() {
	rootCmd.AddCommand(killAllCmd)
}
