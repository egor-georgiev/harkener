package utils

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/google/gopacket/layers"
)

func IntToTCPPort(v int) (layers.TCPPort, error) {
	if v < 0 || v > math.MaxUint16 {
		return 0, fmt.Errorf("ignore port is out of range for a tcp port: %v", v)
	} else {
		return layers.TCPPort(v), nil
	}

}

func ValidateFilePath(path string) (bool, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}
	fileInfo, err := os.Stat(absolutePath)
	if err != nil {
		return false, err
	}
	if fileInfo.IsDir() {
		return false, fmt.Errorf("path is a directory")
	}
	return true, nil
}
