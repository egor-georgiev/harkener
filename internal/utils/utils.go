package utils

import (
	"fmt"
	"math"

	"github.com/google/gopacket/layers"
)

func IntToTCPPort(v int) (layers.TCPPort, error) {
	if v < 0 || v > math.MaxUint16 {
		return 0, fmt.Errorf("ignore port is out of range for a tcp port: %v", v)
	} else {
		return layers.TCPPort(v), nil
	}

}
