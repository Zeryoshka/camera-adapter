package camera

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"

	"errors"

	"github.com/use-go/onvif"
	"github.com/use-go/onvif/media"
	"github.com/use-go/onvif/ptz"
	"github.com/use-go/onvif/xsd"
	onvifTypes "github.com/use-go/onvif/xsd/onvif"
)

type Camera struct {
	ptzParam     *CameraPTZParam
	presetStore  PresetStore
	dev          *onvif.Device
	profileToken onvifTypes.ReferenceToken
}

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
	return &Camera{
		dev:          dev,
		profileToken: profileToken,
		ptzParam: &CameraPTZParam{
			PanSpeed:  0.5,
			TiltSpeed: 0.5,
			ZoomSpeed: 0.5,
		},
		presetStore: make(PresetStore),
	}, nil
}

func getDeviceProfileToken(dev *onvif.Device) (onvifTypes.ReferenceToken, error) {
	resp, err := dev.CallMethod(media.GetProfiles{})
	if err != nil {
		log.Println("Error while getting profiles: ", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)

	type Envelope struct {
		Header struct{}
		Body   struct {
			GetProfilesResponse media.GetProfilesResponse
		}
	}
	reply := Envelope{}
	err = xml.Unmarshal(body, &reply)

	profiles := reply.Body.GetProfilesResponse.Profiles
	if len(profiles) == 0 {
		return "", errors.New("No profiles")
	}

	return profiles[0].Token, nil
}

func (c *Camera) ExecuteCommands(commands ...Command) error {
	for _, command := range commands {
		err := c.ExecuteCommand(command)
		if err != nil {
			return err
		}
	}
	return nil
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

func (c *Camera) executePTZMoveCamera(command *PTZMoveCommand) error {
	newParam := CameraPTZParam{
		PanMove:  command.PanMove,
		TiltMove: command.TiltMove,
		ZoomMove: command.ZoomMove,
	}
	update := c.ptzParam.UpdateMoveParam(&newParam)
	log.Println("Camera params: ", c.ptzParam, "updated: ", update)
	if update {
		err := c.ptzStop()
		if err != nil {
			return err
		}
		return c.ptzContiniousMove()
	}
	return nil
}

func (c *Camera) executePTZPreset(command *PTZPresetCommand) error {
	if command.PresetNumber != 0 {
		if command.SetPreset {
			return c.ptzSetPreset(command.PresetNumber)
		}
		return c.ptzGoToPreset(command.PresetNumber)
	}
	return nil
}

func (c *Camera) ptzStop() error {
	req := ptz.Stop{
		ProfileToken: c.profileToken,
		PanTilt: xsd.Boolean(
			c.ptzParam.PanMove != 0 || c.ptzParam.TiltMove != 0,
		),
		Zoom: xsd.Boolean(c.ptzParam.ZoomMove != 0),
	}
	resp, err := c.dev.CallMethod(req)
	if err != nil {
		log.Fatalln("Error with Stop request: ", err)
		return err
	}
	log.Println("Gotten ", resp.StatusCode, " Stop")
	return nil
}

func (c *Camera) ptzContiniousMove() error {
	req := ptz.ContinuousMove{
		ProfileToken: c.profileToken,
		Velocity: onvifTypes.PTZSpeed{
			PanTilt: onvifTypes.Vector2D{
				X: float64(c.ptzParam.PanMove) * c.ptzParam.PanSpeed,
				Y: float64(c.ptzParam.TiltMove) * c.ptzParam.TiltSpeed,
			},
			Zoom: onvifTypes.Vector1D{
				X: float64(c.ptzParam.ZoomMove) * c.ptzParam.ZoomSpeed,
			},
		},
	}
	resp, err := c.dev.CallMethod(req)
	if err != nil {
		log.Fatalln("Error with ContiniousMove request: ", err)
		return err
	}
	log.Println("Gotten ", resp.StatusCode, " ContiniousMove")
	return nil
}

func (c *Camera) ptzGoToPreset(presetNumber uint) error {
	presetToken, isPresetExist := c.presetStore[presetNumber]
	if !isPresetExist {
		log.Println("No preset with number: ", presetNumber)
		return errors.New(fmt.Sprint("No preset with number: ", presetNumber))
	}
	req := ptz.GotoPreset{
		ProfileToken: c.profileToken,
		PresetToken:  onvifTypes.ReferenceToken(presetToken),
		Speed: onvifTypes.PTZSpeed{
			PanTilt: onvifTypes.Vector2D{
				X: c.ptzParam.PanSpeed,
				Y: c.ptzParam.TiltSpeed,
			},
			Zoom: onvifTypes.Vector1D{
				X: c.ptzParam.ZoomSpeed,
			},
		},
	}
	resp, err := c.dev.CallMethod(req)
	if err != nil {
		log.Println("Error with go to preset \"", presetToken, "\" request: ", err)
		return err
	}
	log.Println("Gotten ", resp.StatusCode, " while go to preset ", presetToken)
	return nil
}

func (c *Camera) ptzSetPreset(presetNumber uint) error {
	presetToken, isPresetExist := c.presetStore[presetNumber]
	if isPresetExist {
		presetToken = ""
	}

	req := ptz.SetPreset{
		ProfileToken: c.profileToken,
		PresetToken:  onvifTypes.ReferenceToken(presetToken),
		PresetName:   xsd.String(fmt.Sprint(presetNumber)),
	}

	resp, err := c.dev.CallMethod(req)
	if err != nil {
		log.Fatalln("Error with set \"", presetToken, "\" preset request: ", err)
		return err
	}
	body, _ := ioutil.ReadAll(resp.Body)

	if !isPresetExist {
		type Envelope struct {
			Header struct{}
			Body   struct {
				SetPresetResponse ptz.SetPresetResponse
			}
		}
		reply := Envelope{}
		xml.Unmarshal(body, &reply)
		c.presetStore[presetNumber] = string(reply.Body.SetPresetResponse.PresetToken)
	}

	log.Println("Gotten ", resp.StatusCode, " while set preset ", presetNumber)
	return nil
}

type CameraPTZParam struct {
	PanMove  int
	TiltMove int
	ZoomMove int

	PanSpeed  float64
	TiltSpeed float64
	ZoomSpeed float64
}

func (p *CameraPTZParam) String() string {
	return fmt.Sprintf("CameraPTZParam{P:%d; T:%d; Z:%d}", p.PanMove, p.TiltMove, p.ZoomMove)
}

func (p *CameraPTZParam) UpdateMoveParam(newP *CameraPTZParam) bool {
	update := false
	if p.PanMove != newP.PanMove || p.TiltMove != newP.TiltMove || p.ZoomMove != newP.ZoomMove {
		update = true
	}

	p.PanMove = newP.PanMove
	p.TiltMove = newP.TiltMove
	p.ZoomMove = newP.ZoomMove

	return update
}

type PresetStore map[uint]string // presetNumber -> token
