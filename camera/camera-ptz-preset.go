package camera

import (
	"errors"
	"fmt"
	"log"

	"github.com/use-go/onvif/ptz"
	"github.com/use-go/onvif/xsd"
	onvifTypes "github.com/use-go/onvif/xsd/onvif"
)

func (c *Camera) executePTZPreset(command *PTZPresetCommand) error {
	if command.PresetNumber != 0 {
		if command.SetPreset {
			return c.ptzSetPreset(command.PresetNumber)
		}
		return c.ptzGoToPreset(command.PresetNumber)
	}
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
		log.Println("Error with GoToPreset \"", presetToken, "\" request: ", err)
		return err
	}
	log.Println("Gotten ", resp.StatusCode, " while GoToPreset (", presetToken, ")")
	return nil
}

func (c *Camera) ptzSetPreset(presetNumber uint) error {
	presetToken := c.presetStore[presetNumber]

	req := ptz.SetPreset{
		ProfileToken: c.profileToken,
		PresetToken:  onvifTypes.ReferenceToken(presetToken),
		PresetName:   xsd.String(fmt.Sprint(presetNumber)),
	}

	resp, err := c.dev.CallMethod(req)
	if err != nil {
		log.Println("Error with SetPreset request: ", err)
		return err
	}

	setPresetResp := ptz.SetPresetResponse{}
	parseSOAPResp(resp, &setPresetResp)
	c.presetStore[presetNumber] = string(setPresetResp.PresetToken)

	log.Println("Gotten ", resp.StatusCode, " while SetPreset(", c.presetStore[presetNumber], ")")
	return nil
}

type PresetStore map[uint]string // presetNumber -> token

func (c *Camera) initPresetStore() error {
	resp, err := c.dev.CallMethod(ptz.GetPresets{ProfileToken: c.profileToken})
	if resp != nil {
		log.Println("Error while getting presets: ", err)
		return err
	}
	presets := ptz.GetPresetsResponse{}
	parseSOAPResp(resp, &presets)
	return nil
}
