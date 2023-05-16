package controller

import (
	"github.com/Zeryoshka/camera-adapter/camera"
	"github.com/Zeryoshka/camera-adapter/confstore"
)

type Controller interface {
	Init() error
	CameraCommandsChan() chan camera.Command
	CamManagerCommandsChan() chan camera.Command
}

func GetController(config *confstore.Config) Controller {
	return NewKeyboardController(config)
}
