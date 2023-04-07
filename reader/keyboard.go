package reader

import (
	"log"
	"strconv"
	"strings"

	"github.com/Zeryoshka/camera-adapter/camera"
	"github.com/Zeryoshka/camera-adapter/confstore"
	"github.com/karalabe/hid"
)

// Use only one of mask and code
type KeyboardProfileKey struct {
	isModifier bool
	mask       byte
	code       byte
}

func (k *KeyboardProfileKey) KeyIsPressed(modyByte byte, pressed map[byte]struct{}) bool {
	if k.isModifier {
		return modyByte&k.mask != 0
	}
	_, isPressed := pressed[k.code]
	return isPressed
}

func NewKeyboardProfileKey(s string) *KeyboardProfileKey {
	keyCode, err := strconv.ParseInt(s, 0, 8)
	if err == nil {
		return &KeyboardProfileKey{
			isModifier: false,
			code:       byte(keyCode),
		}
	}

	if strings.HasPrefix(s, "m") {
		mask, err := strconv.ParseInt(s[1:], 0, 8)
		return &KeyboardProfileKey{
			isModifier: true,
			mask:       byte(mask),
		}
	}

	return nil
}

type ControlProfileKeyboard struct {
	PanLeftKey   *KeyboardProfileKey
	PanRightKey  *KeyboardProfileKey
	TiltUpKey    *KeyboardProfileKey
	TiltDownKey  *KeyboardProfileKey
	ZoomInKey    *KeyboardProfileKey
	ZoomOutKey   *KeyboardProfileKey
	PresetKeys   []*KeyboardProfileKey
	SetPresetKey *KeyboardProfileKey
	UsePresetKey *KeyboardProfileKey
}

type KeyboardReader struct {
	profile *ControlProfileKeyboard
}

func NewKeyboardReader(config confstore.Config) *KeyboardReader {
	configProfile := config.ControlProfile
	profile := &ControlProfileKeyboard{
		PanLeftKey:  NewKeyboardProfileKey(configProfile.PanLeftKey),
		PanRightKey: NewKeyboardProfileKey(configProfile.PanRightKey),
		TiltUpKey:   NewKeyboardProfileKey(configProfile.TiltDownKey),
		TiltDownKey: NewKeyboardProfileKey(configProfile.TiltUpKey),
		ZoomInKey:   NewKeyboardProfileKey(configProfile.ZoomInKey),
		ZoomOutKey:  NewKeyboardProfileKey(configProfile.ZoomOutKey),
	}
	return &KeyboardReader{profile: profile}
}

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

func (r *KeyboardReader) DataToCommands(inputData []byte) ([]camera.Command, []camera.Command) {
	statusKeyByte := inputData[0]
	pressedKeys := make(map[byte]struct{})
	for _, key := range inputData[2:] {
		if key > 3 {
			pressedKeys[key] = struct{}{}
		}
	}
	commands := make([]camera.Command, 0)
	manager_commands := make([]camera.Command, 0)

	panLeftMove := r.profile.PanRightKey.KeyIsPressed(statusKeyByte, pressedKeys)
	panRightMove := r.profile.PanRightKey.KeyIsPressed(statusKeyByte, pressedKeys)
	panMove := 0
	if panLeftMove && !panRightMove {
		panMove = -1
	} else if !panRightMove && panRightMove {
		panMove = +1
	}

	tiltDownMove := r.profile.TiltDownKey.KeyIsPressed(statusKeyByte, pressedKeys)
	tiltUpMove := r.profile.TiltUpKey.KeyIsPressed(statusKeyByte, pressedKeys)
	tiltMove := 0
	if tiltDownMove && !tiltUpMove {
		panMove = -1
	} else if !tiltDownMove && tiltUpMove {
		panMove = +1
	}

	ZoomInMove := r.profile.TiltDownKey.KeyIsPressed(statusKeyByte, pressedKeys)
	ZoomOutMove := r.profile.TiltUpKey.KeyIsPressed(statusKeyByte, pressedKeys)
	zoomMove := 0
	if ZoomInMove && !ZoomOutMove {
		zoomMove = -1
	} else if !ZoomInMove && ZoomOutMove {
		zoomMove = +1
	}

	commands = append(commands, camera.NewPTZMoveCommand(panMove, tiltMove, zoomMove))

	// rightControlPressed := (statusKeyByte & byte(RightControlKeyboardMask)) != 0
	// presetKey, hasPresetCommand := r.pressedOneOf(
	// 	pressedKeys,
	// 	OneNumLockKeyboard, TwoNumLockKeyboard, ThreeNumLockKeyboard, FourNumLockKeyboard, FiveNumLockKeyboard,
	// 	SixNumLockKeyboard, SevenNumLockKeyboard, EightNumLockKeyboard, NineNumLockKeyboard, ZeroNumLockKeyboard,
	// )
	// if hasPresetCommand {
	// 	pressetNumber := uint(presetKey-OneNumLockKeyboard) + 1
	// 	commands = append(commands, camera.NewPTZPresetCommand(rightControlPressed, pressetNumber))
	// }

	// changeDeviceKey, hasChangeDevice := r.pressedOneOf(
	// 	pressedKeys,
	// 	F1Keyboard, F2Keyboard, F3Keyboard, F4Keyboard, F5Keyboard, F6Keyboard, F7Keyboard,
	// 	F8Keyboard, F9Keyboard, F10Keyboard, F11Keyboard, F12Keyboard,
	// )
	// if hasChangeDevice {
	// 	newCameraIndex := int(changeDeviceKey - F1Keyboard)
	// 	manager_commands = append(manager_commands, camera.NewSetDeviceCommand(newCameraIndex))
	// }

	return manager_commands, commands
}
