package main

import (
	"log"
	"net/http"
	"fmt"

	"github.com/girishkoundinya/SAC_Server/router"
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
	fmt.Println("Starting server...");
	router := router.InitRouter()
	log.Fatal(http.ListenAndServe(":3006", router))
}
