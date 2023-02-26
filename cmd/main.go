package main

import (
	"log"

	"github.com/Zeryoshka/camera-adapter/camera"
	"github.com/Zeryoshka/camera-adapter/reader"
	"github.com/use-go/onvif"
)

func main() {
	camera, err := camera.NewCamera(onvif.DeviceParams{
		Xaddr:    "172.18.191.94:80",
		Username: "admin",
		Password: "Supervisor",
	})
	if err != nil {
		log.Fatalln("Can't create camera, cause: ", err)
	}

	reader := reader.GetReader()

	inpCh, err := reader.GetReadChan()
	if err != nil {
		log.Fatalln("Can't get channel: ", err)
	}
	for readedData := range inpCh {
		commands := reader.DataToCommands(readedData)
		log.Println("Got from reader: ", commands, readedData)
		camera.ExecuteCommands(commands...)
	}
}
