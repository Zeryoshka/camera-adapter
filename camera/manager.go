package camera

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Zeryoshka/camera-adapter/confstore"
	"github.com/use-go/onvif"
)

func NewCameraManager(store *confstore.Config) *CameraManager {
	log.Println("Start camera-manager")

	findFirst := false
	startIndex := 0
	cameras := make(map[uint]*Camera)

	for i, cameraConf := range store.Cameras {
		camera, err := NewCamera(
			onvif.DeviceParams{
				Xaddr:    cameraConf.Host,
				Username: cameraConf.Login,
				Password: cameraConf.Password,
				HttpClient: &http.Client{
					Timeout: time.Duration(cameraConf.Timeout) * time.Millisecond,
				},
			},
			cameraConf,
		)
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
	sync.RWMutex
	cameras        map[uint]*Camera
	curCameraIndex uint
}

func (m *CameraManager) GetCamera() *Camera {
	log.Println("Use camera with cameraIndex: ", m.curCameraIndex)
	return m.cameras[m.curCameraIndex]
}

func (m *CameraManager) ExecuteCameraCommand(command Command) error {
	m.RLock()
	defer m.RUnlock()
	return m.GetCamera().ExecuteCommand(command)
}

func (m *CameraManager) ExecuteCommand(command Command) error {
	if command.Type() == SetDeviceCommandType {
		setDeviceCommand := command.(*SetDeviceCommand)
		return m.executeSetDeviceCommand(setDeviceCommand)
	}
	return errors.New(fmt.Sprintln("Unsupported command type for CameraManager: ", command.Type()))
}

func (m *CameraManager) executeSetDeviceCommand(command *SetDeviceCommand) error {
	if m.curCameraIndex == command.CameraIndex {
		log.Printf("Nothing to switch in SetDevice(%d)", m.curCameraIndex)
		return nil
	}
	if cameraPtr, ok := m.cameras[command.CameraIndex]; !ok || cameraPtr == nil {
		return fmt.Errorf("incorrect CameraIndex(%d) in SetDeviceCommand or inactive camera", command.CameraIndex)
	}
	m.Lock()
	defer m.Unlock()
	err := m.cameras[m.curCameraIndex].Stop()
	m.curCameraIndex = command.CameraIndex
	log.Println("Executed command: ", command, " new index updated, cameraIndex: ", m.curCameraIndex)
	if err != nil {
		log.Println("Can't stop previos camera cause: ", err)
	}
	return nil
}
