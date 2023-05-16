package camera

import (
	"errors"
	"fmt"
	"log"

	"github.com/use-go/onvif"
	"github.com/use-go/onvif/media"
	onvifTypes "github.com/use-go/onvif/xsd/onvif"
)

func NewCamera(deviceParams onvif.DeviceParams) (*Camera, error) {
	dev, err := onvif.NewDevice(deviceParams)
	if err != nil {
		log.Println("Error with device creation")
		return nil, err
	}
	profileToken, err := getDeviceProfileToken(dev)
	if err != nil {
		log.Println("Gotten error with getting profile")
		return nil, err
	}
	cam := &Camera{
		dev:          dev,
		profileToken: profileToken,
		ptzParam: &CameraPTZParam{
			PanSpeed:  0.5,
			TiltSpeed: 0.5,
			ZoomSpeed: 0.5,
		},
		presetStore: make(PresetStore),
	}
	if err = cam.initPresetStore(); err != nil {
		return nil, err
	}
	return cam, nil
}

type Camera struct {
	ptzParam     *CameraPTZParam
	presetStore  PresetStore
	dev          *onvif.Device
	profileToken onvifTypes.ReferenceToken
}

func (c *Camera) ExecuteCommand(command Command) error {
	if command.Type() == PTZMoveCommandType {
		ptzMoveCommand := command.(*PTZMoveCommand)
		return c.executePTZMoveCamera(ptzMoveCommand)
	}
	if command.Type() == PTZPresetCommandType {
		ptzPresetCommand := command.(*PTZPresetCommand)
		return c.executePTZPreset(ptzPresetCommand)
	}
	return errors.New(fmt.Sprintln("Unsupported command type for camera: ", command.Type()))
}

func getDeviceProfileToken(dev *onvif.Device) (onvifTypes.ReferenceToken, error) {
	resp, err := dev.CallMethod(media.GetProfiles{})
	if err != nil {
		log.Println("Error while getting profiles: ", err)
		return "", err
	}
	reply := media.GetProfilesResponse{}
	parseSOAPResp(resp, &reply)
	profiles := reply.Profiles
	if len(profiles) == 0 {
		return "", errors.New("NO PROFILES")
	}

	return profiles[len(profiles)-1].Token, nil
}
