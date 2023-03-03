package camera

import (
	"errors"
	"fmt"
	"log"

	"github.com/Zeryoshka/camera-adapter/confstore"
	"github.com/use-go/onvif"
)

type CameraManager struct {
	cameras     map[int]*Camera
	cameraIndex int
}

func NewCameraManager(store *confstore.Store) *CameraManager {
	log.Println("Start creation camera-manager")

	findFirst := false
	startIndex := 0
	cameras := make(map[int]*Camera)
	for index, cameraConf := range store.Cameras {
		camera, err := NewCamera(onvif.DeviceParams{
			Xaddr:    cameraConf.Host,
			Username: cameraConf.Login,
			Password: cameraConf.Password,
		})
		if err != nil {
			log.Println("Can't create camera with host:", cameraConf.Host, ", cause: ", err)
		} else if !findFirst {
			findFirst = true
			startIndex = index
		}
		cameras[index] = camera
	}

	if !findFirst {
		log.Fatalln("Can't open even one camera")
	}
	log.Println("Init ", len(cameras), " cameras for camera-manager")
	return &CameraManager{
		cameras:     cameras,
		cameraIndex: startIndex,
	}
}

func (m *CameraManager) GetCamera() *Camera {
	log.Println("Use camera with cameraIndex: ", m.cameraIndex)
	return m.cameras[m.cameraIndex]
}

func (m *CameraManager) ExecuteCommands(commands ...Command) error {
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
	if cameraPtr, ok := m.cameras[command.CameraIndex]; !ok || cameraPtr == nil {
		return fmt.Errorf("incorrect CameraIndex(%d) in SetDeviceCommand or inactive camera", command.CameraIndex)
	}
	m.cameraIndex = command.CameraIndex
	log.Println("Executed command: ", command, " new index updated, cameraIndex: ", m.cameraIndex)
	return nil
}

type SetDeviceCommand struct {
	CameraIndex int
}

func (c *SetDeviceCommand) String() string {
	return fmt.Sprintf("SetDeviceCommand(index: %d)", c.CameraIndex)
}

func NewSetDeviceCommand(cameraIndex int) *SetDeviceCommand {
	return &SetDeviceCommand{CameraIndex: cameraIndex}
}

func (c *SetDeviceCommand) Type() CommandType {
	return SetDeviceCommandType
}
