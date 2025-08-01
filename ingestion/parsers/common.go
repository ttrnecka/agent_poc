package parsers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

type ExtractFunc func(input string, cfg ExtractorConfig) (any, error)

type ExtractorConfig struct {
	Method  string `yaml:"method"`            // e.g., "match_group", "parse_json"
	Pattern string `yaml:"pattern,omitempty"` // for MatchGroup
	Path    string `yaml:"path,omitempty"`    // for JSONPath, etc.
}

type Config struct {
	Extractors map[string]map[string]ExtractorConfig `yaml:"extractors"`
}

type Endpoint struct {
}

var Extractors = map[string]ExtractFunc{
	"match_group":     MatchGroup,
	"match_group_all": MatchNamedGroupsAll,
	"parse_json":      ParseJSON,
}

func GetCurrentDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to get current file path")
	}
	return filepath.Dir(filename)
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func WatchFile(path string, onChange func()) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	dir := filepath.Dir(absPath)
	fileName := filepath.Base(absPath)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		defer watcher.Close()

		var lastModified time.Time

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// We're only interested in modifications to the target file
				if event.Op&(fsnotify.Write|fsnotify.Create) != 0 && filepath.Base(event.Name) == fileName {
					info, err := os.Stat(absPath)
					if err == nil && info.ModTime() != lastModified {
						lastModified = info.ModTime()
						fmt.Println("Config file changed, reloading...")
						onChange()
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("watch error:", err)
			}
		}
	}()

	return watcher.Add(dir)
}
