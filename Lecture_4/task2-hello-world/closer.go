package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

func runCloser(ctx context.Context) (err error) {
	conn, err := openConn()
	dummy := make([]byte, 2)
	_, err = conn.Read(dummy)
	if err != nil {
		return err
	}
	log.Println("Closing the conn")
	conn.Close()

	log.Println("Opening a new one")
	conn, err = openConn()
	if err != nil {
		return err
	}
	conn.SetReadDeadline(time.Now().Add(time.Second * 19))
	defer conn.Close()

	buff, err := io.ReadAll(conn)
	if err != nil {
		return err
	}

	if !bytes.Equal(buff, []byte(theMessage)) {
		return fmt.Errorf("Incorrect message: expected %s, found %s", theMessage, string(buff))
	}
	return nil
}

var closerCmd = &cobra.Command{
	Use:          "closer",
	SilenceUsage: true,
	Short:        "task2: closer test",
	RunE: func(cmd *cobra.Command, args []string) error {
		var cancel context.CancelFunc
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer cancel()

		return runTest(ctx, runCloser)
	},
}

func init() {
	rootCmd.AddCommand(closerCmd)
}
