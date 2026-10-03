package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

func runSimple(ctx context.Context) (err error) {
	conn, err := openConn()
	if err != nil {
		return err
	}
	conn.SetReadDeadline(time.Now().Add(time.Second * 19))
	defer conn.Close()

	msg, err := io.ReadAll(conn)
	if err != nil {
		return err
	}

	if !bytes.Equal(msg, []byte(theMessage)) {
		return fmt.Errorf("Incorrect message: expected %s, found %s", theMessage, string(msg))
	}
	return nil
}

var simpleCmd = &cobra.Command{
	Use:          "simple",
	SilenceUsage: true,
	Short:        "task2: simple test",
	RunE: func(cmd *cobra.Command, args []string) error {
		var cancel context.CancelFunc
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer cancel()

		return runTest(ctx, runSimple)
	},
}

func init() {
	rootCmd.AddCommand(simpleCmd)
}
