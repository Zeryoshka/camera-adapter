package camera

import "fmt"

func NewSetDeviceCommand(cameraIndex uint) *SetDeviceCommand {
	return &SetDeviceCommand{CameraIndex: cameraIndex}
}

type SetDeviceCommand struct {
	CameraIndex uint
}

func (c *SetDeviceCommand) String() string {
	return fmt.Sprintf("SetDeviceCommand(index: %d)", c.CameraIndex)
}

func (c *SetDeviceCommand) Type() CommandType {
	return SetDeviceCommandType
}
