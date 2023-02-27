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
	RightControlKeyboardMask StatusKeyKeyboardMask = 0b00010000
)

func (r *KeyboardReader) pressedOneOf(pressedKeys map[KeyboardKey]struct{}, wantedKeys ...KeyboardKey) (KeyboardKey, bool) {
	isFirstFound := false
	var pressedKey KeyboardKey
	for _, key := range wantedKeys {
		if _, found := pressedKeys[key]; found {
			if isFirstFound {
				return pressedKey, false
			}
			isFirstFound = true
			pressedKey = key
		}
	}
	return pressedKey, isFirstFound
}

func (r *KeyboardReader) DataToCommands(inputData []byte) []camera.Command {
	statusKeyByte := inputData[0]
	pressedKeys := make(map[KeyboardKey]struct{})
	for _, key := range inputData[2:] {
		if key > 3 {
			pressedKeys[KeyboardKey(key)] = struct{}{}
		}
	}
	commands := make([]camera.Command, 0)

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

	commands = append(commands, camera.NewPTZMoveCommand(panMove, tiltMove, zoomMove))

	rightControlPressed := (statusKeyByte & byte(RightControlKeyboardMask)) != 0
	presetKey, hasPresetCommand := r.pressedOneOf(
		pressedKeys,
		OneNumLockKeyboard, TwoNumLockKeyboard, ThreeNumLockKeyboard, FourNumLockKeyboard, FiveNumLockKeyboard,
		SixNumLockKeyboard, SevenNumLockKeyboard, EightNumLockKeyboard, NineNumLockKeyboard, ZeroNumLockKeyboard,
	)
	if hasPresetCommand {
		pressetNumber := uint(presetKey - OneNumLockKeyboard)
		commands = append(commands, camera.NewPTZPresetCommand(rightControlPressed, pressetNumber))
	}

	return commands
}
