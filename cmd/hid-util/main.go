package main

import (
	"flag"
	"log"

	"github.com/karalabe/hid"
)

func main() {
	dev_num := flag.Int("dev-num", -1, "")
	flag.Parse()
	devices := hid.Enumerate(0, 0)
	log.Printf("Found %d devices", len(devices))
	log.Println("Start reading")
	var use_device *hid.Device = nil
	if *dev_num >= 0 {
		if *dev_num < len(devices) {
			use_device, _ = devices[*dev_num].Open()
		}
	} else {
		for _, device := range devices {
			d, err := device.Open()
			if err == nil {
				use_device = d
				break
			}
		}
	}
	if use_device == nil {
		log.Fatalln("Can't read from devices")
		return
	}

	data := make([]byte, 8)
	for {
		use_device.Read(data)
		log.Println("Gotten:", data)
	}
}
