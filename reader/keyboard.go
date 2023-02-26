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

	_, leftArrowPressed := pressedKeys[LeftArrowKeyboard]
	_, rightArrowPressed := pressedKeys[RightArrowKeyboard]
	panMove := 0
	if leftArrowPressed && !rightArrowPressed {
		panMove = -1
	} else if !leftArrowPressed && rightArrowPressed {
		panMove = +1
	}

	_, upArrowPressed := pressedKeys[UpArrowKeyboard]
	_, downArrowPressed := pressedKeys[DownArrowKeyboard]
	tiltMove := 0
	if downArrowPressed && !upArrowPressed {
		tiltMove = -1
	} else if !downArrowPressed && upArrowPressed {
		tiltMove = +1
	}

	_, minusNumLockKeyboardPressed := pressedKeys[MinusNumLockKeyboard]
	_, plusNumLockKeyboardPressed := pressedKeys[PlusNumLockKeyboard]
	zoomMove := 0
	if minusNumLockKeyboardPressed && !plusNumLockKeyboardPressed {
		zoomMove = -1
	} else if !minusNumLockKeyboardPressed && plusNumLockKeyboardPressed {
		zoomMove = +1
	}

	return []camera.Command{camera.NewPTZMoveCommand(panMove, tiltMove, zoomMove)}
}
