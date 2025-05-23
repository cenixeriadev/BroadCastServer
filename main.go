package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "broadcast-server",
	Short: "A CLI tool for a broadcast server and client",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Start Command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the broadcast server",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		server := NewServer()
		server.Start(port)
		log.Println("Starting Server")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("port", "p", "8080", "Port to listen on")
}

// Connect Command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to the broadcast server",
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("addr")
		connectToServer(addr)
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
	connectCmd.Flags().StringP("addr", "a", "localhost:8080", "Server address")
}
