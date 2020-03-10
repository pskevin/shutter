package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// parseFileCmd represents the beginSnapshot command
var parseFileCmd = &cobra.Command{
	Use:   "ParseFile",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Requires exactly [filePath] positional arguments")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()
			args := strings.Fields(line)

			if len(args) > 0 {
				fmt.Printf("\nRunning `%s`", line)
				switch args[0] {
				case "BeginSnapshot":
					BeginSnapshot(args[1:])
				case "CollectState":
					CollectState(args[1:])
				case "CreateNode":
					CreateNode(args[1:])
				case "KillAll":
					KillAll(args[1:])
				case "PrintSnapshot":
					PrintSnapshot(args[1:])
				case "Receive":
					Receive(args[1:])
				case "ReceiveAll":
					ReceiveAll(args[1:])
				case "Send":
					Send(args[1:])
				default:
					KillAll([]string{})
					return errors.New("Unknown command")
				}
			} else {
				KillAll([]string{})
				return errors.New("Requires non-empty command")
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(parseFileCmd)
}
