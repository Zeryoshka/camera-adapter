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
	PanLeftKey         string   `yaml:"pan_left"`
	PanRightKey        string   `yaml:"pan_right"`
	TiltUpKey          string   `yaml:"tilt_up"`
	TiltDownKey        string   `yaml:"tilt_down"`
	ZoomInKey          string   `yaml:"zoom_in"`
	ZoomOutKey         string   `yaml:"zoom_out"`
	SetPresetKey       string   `yaml:"set_preset"`
	UsePresetKey       string   `yaml:"use_preset"`
	PresetKeys         []string `yaml:"presets"`
	UseChooseCameraKey string   `yaml:"use_choose_camera"`
	ChooseCameraKeys   []string `yaml:"choose_camera"`
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
	// runConfigReloadListener(confpath)
	return &t
}

func runConfigReloadListener(confpath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Can't run fsnotify watcher open: ", err)
	}

	go func() {
		defer watcher.Close()
		for {
			select {
			case event := <-watcher.Events:
				if event.Has(fsnotify.Write) {
					log.Println("modified config, restart MAIN LOOP")
					os.Exit(0)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					log.Fatalln("File watcher retrun error:", err)
				}
				log.Println("File watcher retrun non-critical error:", err)
			}
		}
	}()

	if err = watcher.Add(confpath); err != nil {
		log.Fatal("Can't add file to listener: ", err)
	}
}
