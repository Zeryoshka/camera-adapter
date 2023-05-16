package camera

import (
	"errors"
	"fmt"
	"log"

	"github.com/Zeryoshka/camera-adapter/confstore"
	"github.com/use-go/onvif"
)

func NewCameraManager(store *confstore.Config) *CameraManager {
	log.Println("Start camera-manager")

	findFirst := false
	startIndex := 0
	cameras := make(map[uint]*Camera)

	for i, cameraConf := range store.Cameras {
		camera, err := NewCamera(onvif.DeviceParams{
			Xaddr:    cameraConf.Host,
			Username: cameraConf.Login,
			Password: cameraConf.Password,
		})
		if err != nil {
			log.Println("Can't create camera with host:", cameraConf.Host, ", cause: ", err)
		} else if !findFirst {
			findFirst = true
			startIndex = i
		}
		cameras[uint(i)] = camera
	}

	if !findFirst {
		log.Fatalln("Can't open even one camera")
	}
	log.Println("Init ", len(cameras), " cameras for camera-manager")
	return &CameraManager{
		cameras:        cameras,
		curCameraIndex: uint(startIndex),
	}
}

type CameraManager struct {
	cameras        map[uint]*Camera
	curCameraIndex uint
}

func (m *CameraManager) GetCamera() *Camera {
	log.Println("Use camera with cameraIndex: ", m.curCameraIndex)
	return m.cameras[m.curCameraIndex]
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
	m.curCameraIndex = command.CameraIndex
	log.Println("Executed command: ", command, " new index updated, cameraIndex: ", m.curCameraIndex)
	return nil
}
