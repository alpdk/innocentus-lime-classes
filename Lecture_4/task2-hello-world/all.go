package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	runnerByName = map[string]BodyFunc{
		"simple": runSimple,
		"pacer":  runPacer,
		"closer": runCloser,
		"two":    runTwo,
	}
)

var allCmd = &cobra.Command{
	Use:          "all",
	SilenceUsage: true,
	Short:        "task2: all tests",
	RunE: func(cmd *cobra.Command, args []string) error {
		var cancel context.CancelFunc
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer cancel()

		for _, name := range testsByComplexity {
			runner := runnerByName[name]
			log.Printf("Run: %s\n", name)
			if err := runTest(ctx, runner); err != nil {
				return fmt.Errorf("test %s failed: %v", name, err)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(allCmd)
}
