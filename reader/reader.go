package reader

import (
	"github.com/Zeryoshka/camera-adapter/camera"
	"github.com/Zeryoshka/camera-adapter/confstore"
)

type Reader interface {
	GetReadChan() (<-chan []byte, error)
	DataToCommands([]byte) ([]camera.Command, []camera.Command) // (manager-commands, camera-commands)
}

func GetReader(config *confstore.Config) Reader {
	return NewKeyboardReader(config)
}
