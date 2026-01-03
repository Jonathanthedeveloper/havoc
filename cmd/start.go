/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"time"

	"github.com/Jonathanthedeveloper/havoc.git/internal/proxy"
	"github.com/Jonathanthedeveloper/havoc.git/internal/state"
	"github.com/Jonathanthedeveloper/havoc.git/internal/tui"
	"github.com/spf13/cobra"
)

var (
	port        int
	target      string
	interactive bool
	jitter      time.Duration
	latency     time.Duration
	dropRate    float64
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the Havoc proxy server",
	Long:  `Start a proxy that simulates poor network conditions.`,
	Run: func(cmd *cobra.Command, args []string) {

		state := state.New()

		parsedTarget, err := proxy.ParseTarget(target)
		if err != nil {
			cmd.PrintErrln("Error parsing target:", err)
			return
		}

		state.SetPort(port)
		state.SetTarget(parsedTarget)
		state.SetJitter(jitter)
		state.SetLatency(latency)
		state.SetDropRate(dropRate)

		if interactive {
			tui.Start(state)
			return
		} else {
			proxy.Start(state)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().IntVarP(&port, "port", "p", 3000, "Port to listen on")

	startCmd.Flags().StringVarP(&target, "target", "t", "localhost:8080", "Target server")
	startCmd.MarkFlagRequired("target")

	startCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode")

	startCmd.Flags().DurationVarP(&jitter, "jitter", "j", time.Second*1, "Random latency variance (e.g., 50ms)")

	startCmd.Flags().DurationVarP(&latency, "latency", "l", time.Second*1, "Base latency (e.g., 100ms, 1s)")

	startCmd.Flags().Float64VarP(&dropRate, "drop-rate", "d", 0.0, "Packet drop rate (e.g., 0.1 for 10%)")
}
