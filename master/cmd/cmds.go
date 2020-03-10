package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// beginSnapshotCmd represents the beginSnapshot command
var beginSnapshotCmd = &cobra.Command{
	Use:   "BeginSnapshot",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Requires exactly [nodeID] positional arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) { BeginSnapshot(args) },
}

// collectStateCmd represents the collectState command
var collectStateCmd = &cobra.Command{
	Use:   "CollectState",
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
	Run: func(cmd *cobra.Command, args []string) { CollectState(args) },
}

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
	Run: func(cmd *cobra.Command, args []string) { CreateNode(args) },
}

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
	Run: func(cmd *cobra.Command, args []string) { KillAll(args) },
}

// printSnapshotCmd represents the printSnapshot command
var printSnapshotCmd = &cobra.Command{
	Use:   "PrintSnapshot",
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
	Run: func(cmd *cobra.Command, args []string) { PrintSnapshot(args) },
}

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
	Run: func(cmd *cobra.Command, args []string) { Receive(args) },
}

// receiveAllCmd represents the receiveAll command
var receiveAllCmd = &cobra.Command{
	Use:   "ReceiveAll",
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
	Run: func(cmd *cobra.Command, args []string) { ReceiveAll(args) },
}

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
	Run: func(cmd *cobra.Command, args []string) { Send(args) },
}

func init() {
	rootCmd.AddCommand(beginSnapshotCmd)
	rootCmd.AddCommand(collectStateCmd)
	rootCmd.AddCommand(createNodeCmd)
	rootCmd.AddCommand(killAllCmd)
	rootCmd.AddCommand(printSnapshotCmd)
	rootCmd.AddCommand(receiveCmd)
	rootCmd.AddCommand(receiveAllCmd)
	rootCmd.AddCommand(sendCmd)
}
