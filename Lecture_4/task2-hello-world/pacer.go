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

func runPacer(ctx context.Context) (err error) {
	conn, err := openConn()
	if err != nil {
		return err
	}
	conn.SetReadDeadline(time.Now().Add(time.Second * 19))
	defer conn.Close()

	var msg []byte

	buff := make([]byte, len(theMessage))
	for {
		n, err := conn.Read(buff)
		if n > 3 {
			return fmt.Errorf("The bytes come too fast")
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		msg = append(msg, buff[:n]...)
	}

	if !bytes.Equal(msg, []byte(theMessage)) {
		return fmt.Errorf("Incorrect message: expected %s, found %s", theMessage, string(msg))
	}
	return nil
}

var pacerCmd = &cobra.Command{
	Use:          "pacer",
	SilenceUsage: true,
	Short:        "task2: pacer test",
	RunE: func(cmd *cobra.Command, args []string) error {
		var cancel context.CancelFunc
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer cancel()

		return runTest(ctx, runPacer)
	},
}

func init() {
	rootCmd.AddCommand(pacerCmd)
}
