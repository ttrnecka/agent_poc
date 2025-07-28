package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/ttrnecka/agent_poc/webapi/api"
)

// functions handling refresh process

var refreshMU sync.Mutex

// make refresh blocking and not refresh mutliple times in paraller
func refresh() error {
	refreshMU.Lock()
	defer refreshMU.Unlock()

	requestURL := fmt.Sprintf("http://%s/api/v1/probe", *addr)
	log.Println("Refreshing probes")
	res, err := http.Get(requestURL)
	if err != nil {
		log.Print(fmt.Errorf("probe refresh: %w", err))
		return err
	}
	defer res.Body.Close()

	var probes []api.Probe
	err = json.NewDecoder(res.Body).Decode(&probes)
	if err != nil {
		log.Print(fmt.Errorf("probe body: %w", err))
		return err
	}

	// process probes and make a list of policies that needs downloading
	policies := make(map[string][]string)
	for _, probe := range probes {
		if probe.Collector == *source {
			if policies[probe.Policy] == nil {
				policies[probe.Policy] = []string{probe.Version}
			} else {
				found := false
				for _, version := range policies[probe.Policy] {
					if version == probe.Version {
						found = true
						break
					}
					if !found {
						policies[probe.Policy] = append(policies[probe.Policy], probe.Version)
					}
				}
			}
		}
	}

	// download policies
	// TODO: compare policies with existing policies and delete those no longer needed
	for name, versions := range policies {
		for _, version := range versions {
			file_name := fmt.Sprintf("bin/%s_%s", name, version)
			if _, err := os.Stat(file_name); err != nil {
				err = downloadFile(file_name, fmt.Sprintf("http://%s/api/v1/policy/%s/%s", *addr, name, version))
				if err != nil {
					log.Printf("Error downloading %s: %v", file_name, err)
				}
			} else {
				log.Printf("File %s already exists, skipping download", file_name)
			}
		}
	}

	return nil
}

func downloadFile(filepath string, url string) error {

	fmt.Printf("Downloading %s\n", filepath)
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	err = makeExecutable(filepath)
	return err
}

func makeExecutable(filePath string) error {
	log.Printf("Setting execute permissions on: %s", filePath)
	switch runtime.GOOS {
	case "linux", "darwin":
		// Use chmod to set executable bit (755)
		return os.Chmod(filePath, 0755)
	case "windows":
		// Optionally, ensure .exe extension if it's a binary
		if filepath.Ext(filePath) != ".exe" {
			newPath := filePath + ".exe"
			if err := os.Rename(filePath, newPath); err != nil {
				return fmt.Errorf("rename to .exe failed: %w", err)
			}
		}
		// Windows doesn't need chmod
		return nil
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}
