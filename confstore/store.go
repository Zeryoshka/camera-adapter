package confstore

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
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

type ControlProfileConfig struct {
	PanLeftKey  string `yaml:"pan_left_key"`
	PanRightKey string `yaml:"pan_right_key"`
	TiltUpKey   string `yaml:"tilt_up_key"`
	TiltDownKey string `yaml:"tilt_down_key"`
	ZoomInKey   string `yaml:"zoom_in_key"`
	ZoomOutKey  string `yaml:"zoom_out_key"`
	PresetKeys  string `yaml:"preset_keys"`
}

type Config struct {
	Cameras        []Camera
	ControlProfile ControlProfileConfig `yaml:"control_profile"`
}

func ParseConfigfile(confpath string) *Config {
	filename, _ := filepath.Abs(confpath)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln("Can't open configfile, cause: ", err)
	}

	t := Config{}
	err = yaml.Unmarshal(yamlFile, &t)
	if err != nil {
		log.Fatalln("Can't read configfile, cause: ", err)
	}
	runConfigReloadListener(confpath)
	return &t
}

func runConfigReloadListener(confpath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Can't run fsnotify watcher open: ", err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Has(fsnotify.Write) {
					log.Println("modified config, restart MAIN LOOP")
					os.Exit(0)
				}
			case err := <-watcher.Errors:
				log.Fatalln("File watcher retrun error:", err)
			}
		}
	}()

	if err = watcher.Add(confpath); err != nil {
		log.Fatal("Can't add file to listener: ", err)
	}
}
