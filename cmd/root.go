package cmd

import (
	"harkener/internal/capture"
	"harkener/internal/server"
	"harkener/internal/utils"
	"log"
	"os"
	"sync"

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

		portInfo := make(chan layers.TCPPort)
		var wg sync.WaitGroup

		wg.Add(2)
		go func() {
			defer wg.Done()
			capture.Capture(interfaceName, ignoreTCPPorts, portInfo)
		}()
		go func() {
			defer wg.Done()
			server.Serve(portInfo, bindAddr)
		}()

		wg.Wait()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&interfaceName, "interface", "eth0", "interface to listen on")
	rootCmd.PersistentFlags().StringVar(&bindAddr, "bind", "0.0.0.0:6060", "address to bind to")
	rootCmd.PersistentFlags().IntSliceVar(&ignorePorts, "ignore", []int{}, "ports to ignore")
}
