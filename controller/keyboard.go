package controller

import (
	"log"

	"github.com/Zeryoshka/camera-adapter/camera"
	"github.com/Zeryoshka/camera-adapter/confstore"
	"github.com/karalabe/hid"
)

type KeyboardController struct {
	profile            *ControlProfileKeyboard
	device             *hid.Device
	cameraCommands     chan camera.Command
	camManagerCommands chan camera.Command
}

func NewKeyboardController(config *confstore.Config) *KeyboardController {
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
	return &KeyboardController{profile: profile}
}

func (c *KeyboardController) CameraCommandsChan() <-chan camera.Command {
	return c.cameraCommands
}

func (c *KeyboardController) CamManagerCommandsChan() <-chan camera.Command {
	return c.camManagerCommands
}

func (c *KeyboardController) closeExecution() {
	close(c.camManagerCommands)
	close(c.cameraCommands)
}

func (c *KeyboardController) Init() error {
	devices := hid.Enumerate(0, 0)
	device, err := devices[0].Open()
	if err != nil {
		log.Println("Can't open device, cause: ", err)
		return err
	}
	c.device = device

	go c.startListener()
	return nil
}

func (c *KeyboardController) startListener() {
	data := make([]byte, 8)
	for {
		_, err := c.device.Read(data)
		if err != nil {
			log.Println("Can't read data:", err)
			c.closeExecution()
			return
		}
		c.dataToCommands(data)
	}
}

func (c *KeyboardController) dataToCommands(inputData []byte) {
	statusKeyByte := inputData[0]
	pressedKeys := make(map[byte]struct{})
	for _, key := range inputData[2:] {
		if key > 3 {
			pressedKeys[key] = struct{}{}
		}
	}

	isUsePreset := (c.profile.UsePresetKey == nil) || c.profile.UsePresetKey.KeyIsPressed(statusKeyByte, pressedKeys)
	if isUsePreset {
		isSetPreset := (c.profile.SetPresetKey != nil) && c.profile.SetPresetKey.KeyIsPressed(statusKeyByte, pressedKeys)
		presetKey, keyIndex := c.profile.PresetKeys.PressedOneKey(statusKeyByte, pressedKeys)
		if presetKey != nil {
			c.cameraCommands <- camera.NewPTZPresetCommand(isSetPreset, keyIndex)
			return
		}
	}

	isChooseCamera := (c.profile.UseChooseCameraKey == nil) || c.profile.UseChooseCameraKey.KeyIsPressed(statusKeyByte, pressedKeys)
	if isChooseCamera {
		chooseCameraKey, keyIndex := c.profile.ChooseCameraKeys.PressedOneKey(statusKeyByte, pressedKeys)
		if chooseCameraKey != nil {
			c.camManagerCommands <- camera.NewSetDeviceCommand(keyIndex)
			return
		}
	}

	c.cameraCommands <- camera.NewPTZMoveCommand(
		getPTZMoveInDirection(c.profile.PanLeftKey, c.profile.PanRightKey, statusKeyByte, pressedKeys),
		getPTZMoveInDirection(c.profile.TiltDownKey, c.profile.TiltUpKey, statusKeyByte, pressedKeys),
		getPTZMoveInDirection(c.profile.ZoomOutKey, c.profile.ZoomInKey, statusKeyByte, pressedKeys),
	)
}

func getPTZMoveInDirection(
	negativeKey, positiveKey *KeyboardProfileKey,
	statusKeyByte byte,
	pressedKeys map[byte]struct{},
) int {
	move := 0
	if negativeKey != nil && negativeKey.KeyIsPressed(statusKeyByte, pressedKeys) {
		move -= 1
	}
	if positiveKey != nil && positiveKey.KeyIsPressed(statusKeyByte, pressedKeys) {
		move += 1
	}
	return move
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
	ChooseCameraKeys   *KeyboardProfileKeyGroup
	UseChooseCameraKey *KeyboardProfileKey
}
