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
var tlsCertPath string
var tlsKeyPath string

var rootCmd = &cobra.Command{
	Use:   "harkener",
	Short: "listens to incoming TCP SYN packets, filters and serves them via a websocket server",
	Run: func(cmd *cobra.Command, args []string) {
		ignoreTCPPorts := make(map[layers.TCPPort]struct{})
		for _, port := range ignorePorts {
			castedPort, err := utils.IntToTCPPort(port)
			if err != nil {
				log.Fatalf("failed while casting int to TCP port: %v\n", err)
			}
			ignoreTCPPorts[layers.TCPPort(castedPort)] = struct{}{}
		}
		tlsCertSet := tlsCertPath != ""
		tlsKeySet := tlsKeyPath != ""
		if tlsCertSet != tlsKeySet {
			log.Fatalf("either both tls-cert-path and tls-key-path must be speicifed or none\n")
		}
		if tlsKeySet && tlsCertSet {
			ok, err := utils.ValidateFilePath(tlsKeyPath)
			if !ok {
				log.Fatalf("invalid tls-key-path value: %v\n", err)
			}
			ok, err = utils.ValidateFilePath(tlsCertPath)
			if !ok {
				log.Fatalf("invalid tls-cert-path value: %v\n", err)
			}
		}

		shutdownSignals := make(chan os.Signal, 1)
		signal.Notify(shutdownSignals, syscall.SIGINT, syscall.SIGTERM)

		portInfo := make(chan uint16)

		state, cancel := internal.NewStateWithCancel()

		go capture.Capture(interfaceName, ignoreTCPPorts, portInfo, state)
		go server.Serve(portInfo, bindAddr, state, tlsCertPath, tlsKeyPath)

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
	rootCmd.PersistentFlags().StringVar(&interfaceName, "interface", "", "network interface to capture TCP SYN packets from")
	rootCmd.PersistentFlags().StringVar(&bindAddr, "bind", "", "local host:port for the websocket server to bind to")
	rootCmd.PersistentFlags().IntSliceVar(&ignorePorts, "ignore", []int{}, "ports to ignore")
	rootCmd.MarkPersistentFlagRequired("interface")
	rootCmd.MarkPersistentFlagRequired("bind")
	rootCmd.PersistentFlags().StringVar(&tlsCertPath, "tls-cert-path", "", "path to a TLS certificate")
	rootCmd.PersistentFlags().StringVar(&tlsKeyPath, "tls-key-path", "", "path to a TLS key")
}
