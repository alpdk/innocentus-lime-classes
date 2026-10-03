package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var (
	testsByComplexity = []string{
		"simple",
		"pacer",
		"closer",
		"two",
	}

	theMessage = "Hello, World!"
)

func run(self *cobra.Command, args []string) {
	fmt.Println("The command itself is not a test! Run the tests by name. Here is their list ordered by complexity")
	for idx, name := range testsByComplexity {
		fmt.Printf("%d) %s\n", idx+1, name)
	}
}

type BodyFunc func(ctx context.Context) error

func runTest(ctx context.Context, body BodyFunc) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	testErr := body(ctx)
	if err = ctx.Err(); err != nil {
		return err
	}

	if err = testErr; err != nil {
		log.Printf("test failed: %v\n", err)
		return errors.New("test failed")
	} else {
		log.Println("test passed")
	}
	return nil
}

func openConn() (conn *net.TCPConn, err error) {
	hostStr, portStr, err := net.SplitHostPort(flags.host)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(hostStr)
	if ip == nil {
		return nil, fmt.Errorf("%s is not a valid address", hostStr)
	}

	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("%s is not a valid port number: %v", portStr, err)
	}

	return net.DialTCP("tcp", nil, &net.TCPAddr{IP: ip, Port: int(port)})
}

var rootCmd = &cobra.Command{
	Use:          "client",
	SilenceUsage: true,
	Short:        "task2: hello world",
	Run:          run,
	Long:         "task2: Hello World! Just send \"hello, world\" to every connection. Simple, isn't it? Your deadline is 20 seconds.",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flags.host, "host", "p", "", "server host")
}

var flags struct {
	host string
}

func main() {
	rootCmd.Execute()
}
