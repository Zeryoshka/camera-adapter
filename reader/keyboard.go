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
	Line       string
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
			Line:       s,
		}
	}

	if strings.HasPrefix(s, "m") {
		mask, err := strconv.ParseInt(s[1:], 0, 8)
		if err == nil {
			return &KeyboardProfileKey{
				isModifier: true,
				mask:       byte(mask),
				Line:       s,
			}
		}
	}

	return nil
}

type KeyboardProfileKeyGroup struct {
	keyProfileMap       map[byte]*KeyboardProfileKey
	statusKeyProfileMap map[byte]*KeyboardProfileKey
}

func NewKeyboardProfileKeyGroup(arrS []string) *KeyboardProfileKeyGroup {
	notNil := false
	keyProfileMap := make(map[byte]*KeyboardProfileKey)
	statusKeyProfileMap := make(map[byte]*KeyboardProfileKey)
	for _, s := range arrS {
		keyProfile := NewKeyboardProfileKey(s)
		if keyProfile == nil {
			continue
		} else if keyProfile.isModifier {
			statusKeyProfileMap[keyProfile.mask] = keyProfile
			notNil = true
		} else {
			keyProfileMap[keyProfile.code] = keyProfile
			notNil = true
		}
	}
	if !notNil {
		return nil
	}
	return &KeyboardProfileKeyGroup{
		keyProfileMap:       keyProfileMap,
		statusKeyProfileMap: statusKeyProfileMap,
	}
}

func (g *KeyboardProfileKeyGroup) PressedOneKey(
	statusByte byte, pressed map[byte]struct{},
) *KeyboardProfileKey {
	mask := byte(0b1)
	if statusByte != 0 {
		for {
			keyProfile := g.statusKeyProfileMap[mask]
			if keyProfile != nil && keyProfile.mask == mask {
				return keyProfile
			}
			if mask == 0b10000000 {
				break
			}
			mask <<= 1
		}
	}
	for keyCode := range pressed {
		keyProfile := g.keyProfileMap[keyCode]
		if keyProfile != nil {
			return keyProfile
		}
	}
	return nil
}

type ControlProfileKeyboard struct {
	PanLeftKey         *KeyboardProfileKey
	PanRightKey        *KeyboardProfileKey
	TiltUpKey          *KeyboardProfileKey
	TiltDownKey        *KeyboardProfileKey
	ZoomInKey          *KeyboardProfileKey
	ZoomOutKey         *KeyboardProfileKey
	PresetKeys         *KeyboardProfileKeyGroup
	SetPresetKey       *KeyboardProfileKey
	UsePresetKey       *KeyboardProfileKey
	UseChooseCameraKey *KeyboardProfileKey
	ChooseCameraKeys   *KeyboardProfileKeyGroup
}

type KeyboardReader struct {
	profile *ControlProfileKeyboard
}

func NewKeyboardReader(config *confstore.Config) *KeyboardReader {
	configProfile := config.ControlProfile
	profile := &ControlProfileKeyboard{
		PanLeftKey:         NewKeyboardProfileKey(configProfile.PanLeftKey),
		PanRightKey:        NewKeyboardProfileKey(configProfile.PanRightKey),
		TiltUpKey:          NewKeyboardProfileKey(configProfile.TiltUpKey),
		TiltDownKey:        NewKeyboardProfileKey(configProfile.TiltDownKey),
		ZoomInKey:          NewKeyboardProfileKey(configProfile.ZoomInKey),
		ZoomOutKey:         NewKeyboardProfileKey(configProfile.ZoomOutKey),
		PresetKeys:         NewKeyboardProfileKeyGroup(configProfile.PresetKeys),
		SetPresetKey:       NewKeyboardProfileKey(configProfile.SetPresetKey),
		UsePresetKey:       NewKeyboardProfileKey(configProfile.UsePresetKey),
		ChooseCameraKeys:   NewKeyboardProfileKeyGroup(configProfile.ChooseCameraKeys),
		UseChooseCameraKey: NewKeyboardProfileKey(configProfile.UseChooseCameraKey),
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
			if len(dataChan) < 3 {
				dataChan <- data
			}
		}
	}()

	return dataChan, nil
}

type StatusKeyKeyboardMask byte

const (
	RightControlKeyboardMask StatusKeyKeyboardMask = 0b00010000
)

func (r *KeyboardReader) PressedOneOf(pressedKeys map[KeyboardKey]struct{}, wantedKeys ...KeyboardKey) (KeyboardKey, bool) {
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

	isLockKbMode := false
	isUsePreset := true
	if r.profile.UsePresetKey != nil {
		isUsePreset = r.profile.UsePresetKey.KeyIsPressed(statusKeyByte, pressedKeys)
		isLockKbMode = isUsePreset
	}
	isChooseCamera := !isLockKbMode
	if !isLockKbMode && r.profile.UseChooseCameraKey != nil {
		isChooseCamera = r.profile.UseChooseCameraKey.KeyIsPressed(statusKeyByte, pressedKeys)
	}

	// manage PTZ command
	if !isLockKbMode {
		panLeftMove := false
		if r.profile.PanLeftKey != nil {
			panLeftMove = r.profile.PanLeftKey.KeyIsPressed(statusKeyByte, pressedKeys)
		}
		panRightMove := false
		if r.profile.PanRightKey != nil {
			panRightMove = r.profile.PanRightKey.KeyIsPressed(statusKeyByte, pressedKeys)
		}
		log.Println("DD:", panRightMove)

		panMove := 0
		if panLeftMove && !panRightMove {
			panMove = -1
		} else if !panLeftMove && panRightMove {
			panMove = +1
		}

		tiltDownMove := false
		if r.profile.TiltDownKey != nil {
			tiltDownMove = r.profile.TiltDownKey.KeyIsPressed(statusKeyByte, pressedKeys)
		}
		tiltUpMove := false
		if r.profile.TiltUpKey != nil {
			tiltUpMove = r.profile.TiltUpKey.KeyIsPressed(statusKeyByte, pressedKeys)
		}

		tiltMove := 0
		if tiltDownMove && !tiltUpMove {
			tiltMove = -1
		} else if !tiltDownMove && tiltUpMove {
			tiltMove = +1
		}

		zoomOutMove := false
		if r.profile.ZoomOutKey != nil {
			zoomOutMove = r.profile.ZoomOutKey.KeyIsPressed(statusKeyByte, pressedKeys)
		}
		zoomInMove := false
		if r.profile.ZoomInKey != nil {
			zoomInMove = r.profile.ZoomInKey.KeyIsPressed(statusKeyByte, pressedKeys)
		}

		zoomMove := 0
		if zoomOutMove && !zoomInMove {
			zoomMove = -1
		} else if !zoomOutMove && zoomInMove {
			zoomMove = +1
		}

		commands = append(commands, camera.NewPTZMoveCommand(panMove, tiltMove, zoomMove))
	}

	isSetPreset := false
	if r.profile.SetPresetKey != nil && r.profile.SetPresetKey.KeyIsPressed(statusKeyByte, pressedKeys) {
		isSetPreset = true
	}
	if r.profile.PresetKeys != nil {
		presetKey := r.profile.PresetKeys.PressedOneKey(statusKeyByte, pressedKeys)
		if (presetKey != nil) && (isUsePreset || isSetPreset) {
			pressetNumber := uint(presetKey.code)
			commands = append(commands, camera.NewPTZPresetCommand(isSetPreset, pressetNumber))
		}
	}

	if isChooseCamera {
		chooseCameraKey := r.profile.ChooseCameraKeys.PressedOneKey(statusKeyByte, pressedKeys)
		if chooseCameraKey != nil {
			manager_commands = append(manager_commands, camera.NewSetDeviceCommand(chooseCameraKey.Line))
		}
	}

	// if hasChangeDevice {
	// 	newCameraIndex := int(changeDeviceKey - F1Keyboard)
	// 	manager_commands = append(manager_commands, camera.NewSetDeviceCommand(newCameraIndex))
	// }

	return manager_commands, commands
}
