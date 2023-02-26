package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/karalabe/hid"
)

func main() {
	file, _ := os.Create("devices.txt")
	writer := bufio.NewWriter(file)
	devices := hid.Enumerate(0, 0)
	log.Println(len(devices))
	for _, device := range devices {
		fmt.Println(device.VendorID, device.ProductID, device.Path, device.Serial)
	}
	writer.Flush()
	file.Close()
}
