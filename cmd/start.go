/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	port        int
	target      string
	interactive bool
	jitter      time.Duration
	latency     time.Duration
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the Havoc proxy server",
	Long:  `Start a proxy that simulates poor network conditions.`,
	Run: func(cmd *cobra.Command, args []string) {
		if interactive {
			fmt.Println("Interactive mode enabled")
			return
		} else {
			fmt.Println("Non-interactive mode enabled")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().IntVarP(&port, "port", "p", 3000, "Port to listen on")

	startCmd.Flags().StringVarP(&target, "target", "t", "localhost:8080", "Target server")
	startCmd.MarkFlagRequired("target")

	startCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode")

	startCmd.Flags().DurationVarP(&jitter, "jitter", "j", time.Second*5, "Random latency variance (e.g., 50ms)")

	startCmd.Flags().DurationVarP(&latency, "latency", "l", time.Second*5, "Base latency (e.g., 100ms, 1s)")
}
