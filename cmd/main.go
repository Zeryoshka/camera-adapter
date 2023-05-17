package main

import (
	"flag"
	"log"

	"github.com/Zeryoshka/camera-adapter/camera"
	"github.com/Zeryoshka/camera-adapter/confstore"
	"github.com/Zeryoshka/camera-adapter/controller"
)

func main() {
	confpath := flag.String("conf", "", "path to local yaml config")
	flag.Parse()
	config := confstore.ParseConfigfile(*confpath)

	manager := camera.NewCameraManager(config)

	controller := controller.GetController(config)
	err := controller.Init()
	if err != nil {
		log.Fatalln("Can't init controller: ", err)
	}
	log.Println("Start Upravlaytor")
	camManagerChan := controller.CamManagerCommandsChan()
	camerChan := controller.CameraCommandsChan()
	for {
		select {
		case command := <-camManagerChan:
			manager.ExecuteCommand(command)
		case command := <-camerChan:
			manager.ExecuteCameraCommand(command)
		}
	}
}
