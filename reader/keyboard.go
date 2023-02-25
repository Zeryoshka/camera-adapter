package reader

import (
	"log"

	"github.com/Zeryoshka/camera-adapter/camera"
	"github.com/karalabe/hid"
)

type KeyboardReader struct{}

func (r *KeyboardReader) GetReadChan() (<-chan []byte, error) {
	devices := hid.Enumerate(0, 0)
	device, err := devices[0].Open()

	if err != nil {
		log.Println("Can't open device, cause: ", err)
		return nil, err
	}

	dataChan := make(chan []byte)
	go func() {
		data := make([]byte, 8)
		for {
			_, err := device.Read(data)
			if err != nil {
				log.Println("Can't read data:", err)
				close(dataChan)
				return
			}
			dataChan <- data
		}
	}()

	return dataChan, nil
}

type KeyboardKey byte

const (
	RightArrowKeyboard   KeyboardKey = 79
	LeftArrowKeyboard    KeyboardKey = 80
	DownArrowKeyboard    KeyboardKey = 81
	UpArrowKeyboard      KeyboardKey = 82
	MinusNumLockKeyboard KeyboardKey = 86
	PlusNumLockKeyboard  KeyboardKey = 87
	OneNumLockKeyboard   KeyboardKey = 89
	TwoNumLockKeyboard   KeyboardKey = 90
	ThreeNumLockKeyboard KeyboardKey = 91
	FourNumLockKeyboard  KeyboardKey = 92
	FiveNumLockKeyboard  KeyboardKey = 93
	SixNumLockKeyboard   KeyboardKey = 94
	SevenNumLockKeyboard KeyboardKey = 95
	EightNumLockKeyboard KeyboardKey = 96
	NineNumLockKeyboard  KeyboardKey = 97
	ZeroNumLockKeyboard  KeyboardKey = 98
)

type StatusKeyKeyboardMask byte

const (
	RightControl StatusKeyKeyboardMask = 0b00000100
)

func (r *KeyboardReader) DataToCommands(inputData []byte) []camera.Command {
	pressedKeys := make(map[KeyboardKey]struct{})
	for _, key := range inputData[2:] {
		if key > 3 {
			pressedKeys[KeyboardKey(key)] = struct{}{}
		}
	}

	var panCommand camera.Command
	_, leftArrowPressed := pressedKeys[LeftArrowKeyboard]
	_, rightArrowPressed := pressedKeys[RightArrowKeyboard]
	if leftArrowPressed && !rightArrowPressed {
		panCommand.Action = camera.PanLeftAction
	} else if rightArrowPressed && !leftArrowPressed {
		panCommand.Action = camera.PanRightAction
	} else {
		panCommand.Action = camera.PanStopAction
	}

	var tiltCommand camera.Command
	_, upArrowPressed := pressedKeys[UpArrowKeyboard]
	_, downArrowPressed := pressedKeys[DownArrowKeyboard]
	if upArrowPressed && !downArrowPressed {
		tiltCommand.Action = camera.TiltUpAction
	} else if downArrowPressed && !upArrowPressed {
		tiltCommand.Action = camera.TiltDownAction
	} else {
		tiltCommand.Action = camera.TiltStopAction
	}
	return []camera.Command{panCommand, tiltCommand}
}
