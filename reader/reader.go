package reader

import (
	"github.com/Zeryoshka/camera-adapter/camera"
)

type Reader interface {
	GetReadChan() (<-chan []byte, error)
	DataToCommands([]byte) ([]camera.Command, []camera.Command) // (manager-commands, camera-commands)
}

func GetReader() Reader {
	return &KeyboardReader{}
}
