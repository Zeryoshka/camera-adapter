package main

import (
	"flag"
	"log"

	"github.com/Zeryoshka/camera-adapter/camera"
	"github.com/Zeryoshka/camera-adapter/confstore"
	"github.com/Zeryoshka/camera-adapter/reader"
)

func main() {
	confpath := flag.String("conf", "", "path to local yaml config")
	flag.Parse()
	store := confstore.NewFileStore(*confpath)
	manager := camera.NewCameraManager(store)

	reader := reader.GetReader()

	inpCh, err := reader.GetReadChan()
	if err != nil {
		log.Fatalln("Can't get channel: ", err)
	}
	for readedData := range inpCh {
		manager_commands, cam_commands := reader.DataToCommands(readedData)
		log.Println("Got from reader: ", manager_commands, cam_commands, readedData)
		manager.ExecuteCommands(manager_commands...)
		cam := manager.GetCamera()
		cam.ExecuteCommands(cam_commands...)
	}
}
