package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

func downloadFile(success chan<- error, conn *net.TCPConn, fileHash []byte) {
	netHash, err := hashReader(conn)
	if err != nil {
		success <- err
	}

	if !bytes.Equal(fileHash, netHash) {
		success <- fmt.Errorf("Hash mismatched: expected %x, got %x", fileHash, netHash)
	}
	success <- nil
}

func runTest(ctx context.Context, conn *net.TCPConn, fileHash []byte) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	success := make(chan error)
	go downloadFile(success, conn, fileHash)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err = <-success:
		return err
	}
}

func hashFile(path string) (hash []byte, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return hashReader(file)
}

func hashReader(r io.Reader) (hash []byte, err error) {
	h := sha256.New()
	buff := make([]byte, 4096)
	for {
		n, err := r.Read(buff)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		_, err = h.Write(buff[:n])
		if err != nil {
			return nil, err
		}
	}
	return h.Sum(nil), nil
}

func runE(self *cobra.Command, args []string) (err error) {
	var cancel context.CancelFunc
	ctx, cancel := signal.NotifyContext(self.Context(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	hostStr, portStr, err := net.SplitHostPort(flags.host)
	if err != nil {
		return err
	}

	ip := net.ParseIP(hostStr)
	if ip == nil {
		return fmt.Errorf("%s is not a valid address", hostStr)
	}

	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return fmt.Errorf("%s is not a valid port number: %v", portStr, err)
	}

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: ip, Port: int(port)})
	if err != nil {
		return err
	}
	defer conn.Close()

	hash, err := hashFile(flags.filePath)
	if err != nil {
		return err
	}

	if err = runTest(ctx, conn, hash); err != nil {
		log.Printf("test failed: %v\n", err)
	} else {
		log.Println("test passed")
	}

	return nil
}

var rootCmd = &cobra.Command{
	Use:          "client",
	SilenceUsage: true,
	Short:        "task1: file delivery",
	RunE:         runE,
	Long:         "task1: file delivery. You must write a server using tcp, that can deliver a file. The file must be delivered in 5 seconds.",
}

func init() {
	rootCmd.Flags().StringVarP(&flags.host, "host", "p", "", "server host")
	rootCmd.Flags().StringVarP(&flags.filePath, "file", "i", "", "the file your server will be delivering. I need to compare the integrity")
}

var flags struct {
	host     string
	filePath string
}

func main() {
	rootCmd.Execute()
}
