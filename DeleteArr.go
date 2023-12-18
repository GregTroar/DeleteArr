package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/yaml.v2"
)

type GotifyMessage struct {
	Title    string `json:"title"`
	Priority string `json:"priority"`
	Message  string `json:"message"`
}

type MediaFiles struct {
	EventType    string
	SourcePath   string
	SourceFolder string
	FileName     string
	InFolder     bool
	Arr          string
}

type Config struct {
	Gotify struct {
		Enabled   bool   `yaml:"enabled"`
		ServerURL string `yaml:"server_url"`
		Token     string `yaml:"token"`
	} `yaml:"gotify"`
	General struct {
		RootFolders []string `yaml:"root_folders"`
	} `yaml:"general"`
}

func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func (m *MediaFiles) SendGotify(message string, arr string, cfg *Config) {
	if cfg.Gotify.Enabled {
		http.PostForm(cfg.Gotify.ServerURL+"/message?token="+cfg.Gotify.Token,
			url.Values{"message": {message}, "title": {"Deleting Media from " + arr}, "priority": {"10"}})
	}
}

func (m *MediaFiles) IsInFolder(cfg *Config) {
	folderList := cfg.General.RootFolders
	// []string{"Movies", "4K-Movies", "Series", "4K-Series", "Kids", "Animes"}

	SplitFolder := strings.Split(m.SourcePath, "/")
	LastFolder := SplitFolder[len(SplitFolder)-2]
	m.FileName = SplitFolder[len(SplitFolder)-1]

	log.Printf("Found last folder to be: %v", LastFolder)

	ContainsFolder := slices.Contains(folderList, LastFolder)

	if !ContainsFolder {
		m.InFolder = true
		log.Printf("Movie %v is in a Folder\n", m.FileName)
	} else {
		m.InFolder = false
		log.Printf("Movie %v is not in a Folder\n", m.FileName)
	}
}

func main() {

	// os.Setenv("radarr_moviefile_sourcepath", "/mnt/Multimedia/Download/PostProcess/Movies/Butchers.Crossing.2023.MULTi.1080p.WEB.x264-FW.mkv")
	// os.Setenv("radarr_moviefile_sourcefolder", "/mnt/Multimedia/Download/PostProcess/Movies")
	// os.Setenv("radarr_eventtype", "Download")
	// os.Setenv("EventType", "Test")

	if os.Getenv("radarr_eventtype") == "Test" || os.Getenv("sonarr_eventtype") == "Test" {
		log.Println("Radarr/Sonarr is testing the script and it works")
		os.Exit(0)
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	configPath := filepath.Join(exPath, "config.yml")

	cfg, err := NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(exPath+"/log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	RadarrEventType := os.Getenv("radarr_eventtype")
	SonarrEventType := os.Getenv("sonarr_eventtype")

	m := &MediaFiles{}

	if RadarrEventType != "" {
		m.EventType = os.Getenv("radarr_eventtype")
		m.SourcePath = os.Getenv("radarr_moviefile_sourcepath")
		log.Printf("The Source Path is: %v", m.SourcePath)
		m.SourceFolder = os.Getenv("radarr_moviefile_sourcefolder")
		log.Printf("The Source Folder is: %v", m.SourceFolder)
		m.Arr = "Radarr"

	}

	if SonarrEventType != "" {
		m.EventType = os.Getenv("sonarr_eventtype")
		m.SourcePath = os.Getenv("sonarr_episodefile_sourcepath")
		log.Printf("The Source Path is: %v", m.SourcePath)
		m.SourceFolder = os.Getenv("radarr_moviefile_sourcefolder")
		log.Printf("The Source Folder is: %v", m.SourceFolder)
		m.Arr = "Sonarr"
	}

	m.IsInFolder(cfg)

	if m.InFolder {

		f, err := os.Open(m.SourceFolder)
		if err != nil {
			log.Println(err)
			return
		}
		files, err := f.Readdir(0)
		if err != nil {
			log.Println(err)
			return
		}

		mkvCount := 0

		for _, v := range files {
			log.Println("Found file: " + v.Name())
			if filepath.Ext(v.Name()) == ".mkv" {
				mkvCount += 1
			} else {
				os.Remove(m.SourceFolder + "/" + v.Name())
				log.Printf("Deleting non MKV file: %v", m.SourceFolder+"/"+v.Name())
			}
		}

		if mkvCount > 1 {
			log.Printf("Found %v MKV files in the folder, deleting only %v", mkvCount, m.SourcePath)
			os.Remove(m.SourcePath)
		} else {
			log.Printf("Found only one MKV files in the folder, deleting the folder %v", m.SourceFolder)
			os.RemoveAll(m.SourceFolder)
		}

		// if not in folder just delete the file
	} else {
		os.RemoveAll(m.SourcePath)
		log.Printf("Deleting the file %v", m.SourcePath)
	}

	m.SendGotify("Deleting Source Path"+m.SourcePath, m.Arr, cfg)

}
