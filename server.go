package main

import (
	"github.com/spf13/cobra"
)

func main() {

	var cmdStartServer = &cobra.Command{
		Use:   "start",
		Short: "Command to start server",
		Long:  `start is for starting the server`,
		Run: func(cmd *cobra.Command, args []string) {
			startServer()
		},
	}

	var rootCmd = &cobra.Command{Use: "server"}
	rootCmd.AddCommand(cmdStartServer)
	rootCmd.Execute()
}

func startServer() {

}
