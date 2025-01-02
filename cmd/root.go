package cmd

import (
	"fmt"
	"listener/internal"
	"log"
	"math"
	"os"

	"github.com/google/gopacket/layers"
	"github.com/spf13/cobra"
)

var interfaceName string
var ignorePorts []int

// TODO move to utils
func intToTCPPort(v int) (layers.TCPPort, error) {
	if v < 0 || v > math.MaxUint16 {
		return 0, fmt.Errorf("ignore port is out of range for a tcp port: %v", v)
	} else {
		return layers.TCPPort(v), nil
	}

}

var rootCmd = &cobra.Command{
	Use: "listener",
	//Short: "Listens for TCP SYN packets",
	Run: func(cmd *cobra.Command, args []string) {
		ignoreTCPPorts := make(map[layers.TCPPort]struct{})
		for _, port := range ignorePorts {
			castedPort, err := intToTCPPort(port)
			if err != nil {
				log.Fatalf("failed while casting int to TCP port: %v\n", err)
			}
			ignoreTCPPorts[layers.TCPPort(castedPort)] = struct{}{}
		}

		portInfo := make(chan layers.TCPPort)
		go internal.Listen(interfaceName, ignoreTCPPorts, portInfo)
		for {
			port := <-portInfo
			log.Printf("%d", port)
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
	rootCmd.PersistentFlags().StringVar(&interfaceName, "interface", "eth0", "interface to listen on")
	rootCmd.PersistentFlags().IntSliceVar(&ignorePorts, "ignore", []int{22, 80, 443}, "ports to ignore")
}
