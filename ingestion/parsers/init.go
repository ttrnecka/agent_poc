package parsers

import (
	"fmt"
	"log"
	"path/filepath"
)

var Parsers map[string]*Config

func init() {
	Parsers = make(map[string]*Config)
	err := setupConfig("brocade_cli")
	if err != nil {
		log.Fatal(err)
	}
}

func setupConfig(policy string) error {
	file := filepath.Join("./config", fmt.Sprintf("%s.yaml", policy))
	cfg, err := LoadConfig(file)
	if err != nil {
		return err
	}
	Parsers[policy] = cfg

	err = WatchFile(file, func() {
		cfg, err := LoadConfig(file)
		if err != nil {
			fmt.Printf("Error reloading config: %v\n", err)
			return
		}
		log.Printf("Reloaded config for %s policy", policy)
		Parsers[policy] = cfg
	})
	if err != nil {
		return err
	}
	return nil
}
