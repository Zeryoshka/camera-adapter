package camera

import (
	"errors"
	"fmt"
)

type CameraManager struct {
	cameras     []*Camera
	cameraIndex uint
}

func NewCameraManager() *CameraManager {
	return &CameraManager{}
}

func (m *CameraManager) GetCamera() *Camera {
	return m.cameras[m.cameraIndex]
}

func (m *CameraManager) ExecuteManagerCommands(commands ...Command) error {
	for _, command := range commands {
		err := m.ExecuteCommand(command)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *CameraManager) ExecuteCommand(command Command) error {
	if command.Type() == SetDeviceCommandType {
		setDeviceCommand := command.(*SetDeviceCommand)
		return m.executeSetDeviceCommand(setDeviceCommand)
	}
	return errors.New(fmt.Sprintln("Unsupported command type for CameraManager: ", command.Type()))
}

func (m *CameraManager) executeSetDeviceCommand(command *SetDeviceCommand) error {
	if command.СameraIndex >= 10 {
		return fmt.Errorf("incorrect CameraIndex(%d) in SetDeviceCommand", command.СameraIndex)
	}
	m.cameraIndex = command.СameraIndex
	return nil
}

type SetDeviceCommand struct {
	СameraIndex uint
}

func NewSetDeviceCommand() *SetDeviceCommand {
	return &SetDeviceCommand{}
}

func (c *SetDeviceCommand) Type() CommandType {
	return SetDeviceCommandType
}

func (c *SetDeviceCommand) String() string {
	return "SetDeviceCommand"
}
