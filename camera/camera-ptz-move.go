package camera

import (
	"fmt"
	"log"

	"github.com/use-go/onvif/ptz"
	"github.com/use-go/onvif/xsd"
	onvifTypes "github.com/use-go/onvif/xsd/onvif"
)

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
		log.Println("Error with Stop request: ", err)
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
		Timeout: xsd.Duration("").NewDateTime("0", "0", "0", "0", "0", "1"),
	}

	resp, err := c.dev.CallMethod(req)
	if err != nil {
		log.Fatalln("Error with ContiniousMove request: ", err)
		return err
	}
	log.Println("Gotten ", resp.StatusCode, " ContiniousMove")
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
	updated := (p.PanMove != newP.PanMove) || (p.TiltMove != newP.TiltMove) || (p.ZoomMove != newP.ZoomMove)

	p.PanMove = newP.PanMove
	p.TiltMove = newP.TiltMove
	p.ZoomMove = newP.ZoomMove
	return updated
}
