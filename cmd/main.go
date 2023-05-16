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
	cam := manager.GetCamera()

	controller := controller.GetController(config)
	err := controller.Init()
	if err != nil {
		log.Fatalln("Can't init controller: ", err)
	}

	camManagerChan := controller.CamManagerCommandsChan()
	camerChan := controller.CameraCommandsChan()
	for {
		select {
		case command := <-camManagerChan:
			manager.ExecuteCommand(command)
		case command := <-camerChan:
			cam.ExecuteCommand(command)
		}
	}
}
