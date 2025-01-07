package cmd

import (
	"harkener/internal"
	"harkener/internal/capture"
	"harkener/internal/server"
	"harkener/internal/utils"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/gopacket/layers"
	"github.com/spf13/cobra"
)

var interfaceName string
var bindAddr string
var ignorePorts []int

var rootCmd = &cobra.Command{
	Use:   "harkener",
	Short: "listens to incoming TCP SYN packets, filters and serves them via UDP",
	Run: func(cmd *cobra.Command, args []string) {
		ignoreTCPPorts := make(map[layers.TCPPort]struct{})
		for _, port := range ignorePorts {
			castedPort, err := utils.IntToTCPPort(port)
			if err != nil {
				log.Fatalf("failed while casting int to TCP port: %v\n", err)
			}
			ignoreTCPPorts[layers.TCPPort(castedPort)] = struct{}{}
		}

		shutdownSignals := make(chan os.Signal)
		signal.Notify(shutdownSignals, syscall.SIGINT, syscall.SIGTERM)

		portInfo := make(chan layers.TCPPort)

		state, cancel := internal.NewStateWithCancel()

		go capture.Capture(interfaceName, ignoreTCPPorts, portInfo, state)
		go server.Serve(portInfo, bindAddr, state)

		select {
		case sig := <-shutdownSignals:
			log.Printf("got signal: %v, shutting down gracefully\n", sig)
			cancel()
			close(portInfo)
			log.Printf("bye!\n")
		case err := <-state.Errors:
			log.Printf("got error during execution: %v, shutting down\n", err)
			cancel()
			close(portInfo)
			log.Fatalf("bye!\n")
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&interfaceName, "interface", "", "interface to listen on")
	rootCmd.PersistentFlags().StringVar(&bindAddr, "bind", "", "address to bind to")
	rootCmd.PersistentFlags().IntSliceVar(&ignorePorts, "ignore", []int{}, "ports to ignore")
	rootCmd.MarkPersistentFlagRequired("interface")
	rootCmd.MarkPersistentFlagRequired("bind")
}
