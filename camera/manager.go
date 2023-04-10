package camera

import (
	"errors"
	"fmt"
	"log"

	"github.com/Zeryoshka/camera-adapter/confstore"
	"github.com/use-go/onvif"
)

type CameraManager struct {
	cameras   map[string]*Camera
	cameraKey string
}

func NewCameraManager(store *confstore.Config) *CameraManager {
	log.Println("Start creation camera-manager")

	findFirst := false
	startKey := ""
	cameras := make(map[string]*Camera)
	for _, cameraConf := range store.Cameras {
		camera, err := NewCamera(onvif.DeviceParams{
			Xaddr:    cameraConf.Host,
			Username: cameraConf.Login,
			Password: cameraConf.Password,
		})
		if err != nil {
			log.Println("Can't create camera with host:", cameraConf.Host, ", cause: ", err)
		} else if !findFirst {
			findFirst = true
			startKey = cameraConf.CameraKey
		}
		cameras[cameraConf.CameraKey] = camera
	}

	if !findFirst {
		log.Fatalln("Can't open even one camera")
	}
	log.Println("Init ", len(cameras), " cameras for camera-manager")
	return &CameraManager{
		cameras:   cameras,
		cameraKey: startKey,
	}
}

func (m *CameraManager) GetCamera() *Camera {
	log.Println("Use camera with cameraIndex: ", m.cameraKey)
	return m.cameras[m.cameraKey]
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
	if cameraPtr, ok := m.cameras[command.CameraKey]; !ok || cameraPtr == nil {
		return fmt.Errorf("incorrect CameraIndex(%s) in SetDeviceCommand or inactive camera", command.CameraKey)
	}
	m.cameraKey = command.CameraKey
	log.Println("Executed command: ", command, " new index updated, cameraIndex: ", m.cameraKey)
	return nil
}

type SetDeviceCommand struct {
	CameraKey string
}

func (c *SetDeviceCommand) String() string {
	return fmt.Sprintf("SetDeviceCommand(key: %s)", c.CameraKey)
}

func NewSetDeviceCommand(cameraKey string) *SetDeviceCommand {
	return &SetDeviceCommand{CameraKey: cameraKey}
}

func (c *SetDeviceCommand) Type() CommandType {
	return SetDeviceCommandType
}
