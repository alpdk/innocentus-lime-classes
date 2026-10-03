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

func runTwo(ctx context.Context) (err error) {
	conn1, err := openConn()
	if err != nil {
		return err
	}
	go func() {
		time.Sleep(5 * time.Minute) // Sorry, I don't know how to use it
		conn1.Close()
	}()

	time.Sleep(1 * time.Second)

	conn2, err := openConn()
	if err != nil {
		return err
	}
	defer conn2.Close()
	conn2.SetReadDeadline(time.Now().Add(time.Second * 19))

	msg, err := io.ReadAll(conn2)
	if err != nil {
		return err
	}

	if !bytes.Equal(msg, []byte(theMessage)) {
		return fmt.Errorf("Incorrect message: expected %s, found %s", theMessage, string(msg))
	}
	return nil
}

var twoCmd = &cobra.Command{
	Use:          "two",
	SilenceUsage: true,
	Short:        "task2: two test",
	RunE: func(cmd *cobra.Command, args []string) error {
		var cancel context.CancelFunc
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer cancel()

		return runTest(ctx, runTwo)
	},
}

func init() {
	rootCmd.AddCommand(twoCmd)
}
