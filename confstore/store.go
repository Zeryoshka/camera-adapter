package confstore

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Store struct {
	Cameras []Camera
}

func NewFileStore(confpath string) *Store {
	conf := ParseConfigfile(confpath)
	return &Store{
		Cameras: conf.Cameras,
	}
}

type Camera struct {
	Host     string
	Login    string
	Password string
}

type ConfigFileData struct {
	Cameras []Camera
}

func ParseConfigfile(confpath string) *ConfigFileData {
	filename, _ := filepath.Abs(confpath)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln("Can't open configfile, cause: ", err)
	}

	t := ConfigFileData{}
	err = yaml.Unmarshal(yamlFile, &t)
	if err != nil {
		log.Fatalln("Can't read configfile, cause: ", err)
	}
	return &t
}
