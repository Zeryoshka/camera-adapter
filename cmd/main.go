package main

import (
	"bufio"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Zeryoshka/camera-adapter/reader"
	"github.com/karalabe/hid"
	"github.com/use-go/onvif"
	"github.com/use-go/onvif/media"
	"github.com/use-go/onvif/ptz"
	"github.com/use-go/onvif/xsd"
	onvifTypes "github.com/use-go/onvif/xsd/onvif"
)

func getProfileToken(dev *onvif.Device) string {
	type Envelope struct {
		Header struct{}
		Body   struct {
			GetProfilesResponse media.GetProfilesResponse
		}
	}
	resp, err := dev.CallMethod(media.GetProfiles{})
	if err != nil {
		log.Fatalln("Error while getting profiles: ", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))

	reply := Envelope{}
	err = xml.Unmarshal(body, &reply)
	if err != nil {
		log.Fatalln("Error while getting profiles: ", err)
	}
	profilesResponse := reply.Body.GetProfilesResponse
	log.Println(len(profilesResponse.Profiles))
	for _, profile := range profilesResponse.Profiles {
		log.Println("Profile: ", profile.Name, profile.Token)
	}
	return ""
}

func checkCamera(ctx context.Context) {
	dev, err := onvif.NewDevice(onvif.DeviceParams{
		Xaddr:    "172.18.191.94:80",
		Username: "admin",
		Password: "Supervisor",
	})
	if err != nil {
		log.Fatalln("Can't connect to camera: ", err)
	}

	getProfileToken(dev)

	mvReq := ptz.ContinuousMove{
		Velocity: onvifTypes.PTZSpeed{
			PanTilt: onvifTypes.Vector2D{
				X: 1.0,
				Y: 1.0,
			},
			Zoom: onvifTypes.Vector1D{X: 0.0},
		},
		Timeout:      xsd.Duration("1000"),
		ProfileToken: onvifTypes.ReferenceToken("PROFILE_001"),
	}
	resp, err := dev.CallMethod(mvReq)
	if err != nil {
		log.Fatalln("Error with request: ", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("Onvif:", resp.StatusCode, string(body))
}

func main() {
	ctx := context.Background()
	checkCamera(ctx)
	reader := reader.GetReader()

	inpCh, err := reader.GetReadChan()
	if err != nil {
		log.Fatalln("Can't get channel: ", err)
	}
	for readedData := range inpCh {
		commands := reader.DataToCommands(readedData)
		log.Println(commands, readedData)
	}

	file, _ := os.Create("devices.txt")
	writer := bufio.NewWriter(file)
	devices := hid.Enumerate(0, 0)
	log.Println(len(devices))
	for _, device := range devices {
		writer.WriteString(fmt.Sprintln(
			device.VendorID, device.ProductID, device.Path, device.Serial,
		))
	}
	writer.Flush()
	file.Close()
}
