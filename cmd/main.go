package main

import (
	"flag"
	"log"

	"github.com/Zeryoshka/camera-adapter/camera"
	"github.com/Zeryoshka/camera-adapter/confstore"
	"github.com/Zeryoshka/camera-adapter/reader"
	"github.com/use-go/onvif"
)

func main() {
	confpath := flag.String("conf", "", "path to local yaml config")
	flag.Parse()
	store := confstore.NewFileStore(*confpath)
	camerConf := store.Cameras[0]
	camera, err := camera.NewCamera(onvif.DeviceParams{
		Xaddr:    camerConf.Host,     //"172.18.191.94:80",
		Username: camerConf.Login,    //"admin",
		Password: camerConf.Password, //"Supervisor",
	})
	if err != nil {
		log.Fatalln("Can't create camera, cause: ", err)
	}

	reader := reader.GetReader()

	inpCh, err := reader.GetReadChan()
	if err != nil {
		log.Fatalln("Can't get channel: ", err)
	}
	for readedData := range inpCh {
		commands := reader.DataToCommands(readedData)
		log.Println("Got from reader: ", commands, readedData)
		camera.ExecuteCommands(commands...)
	}
}
