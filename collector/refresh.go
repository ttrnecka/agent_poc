package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
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

	jar, err := cookiejar.New(nil)
	if err != nil {
		logger.Error().Err(err).Msg("")
		return err
	}
	client := &http.Client{
		Jar: jar,
	}

	form := url.Values{}
	form.Set("username", "test")
	form.Set("password", "test")

	resp, err := client.PostForm(fmt.Sprintf("http://%s/login", *addr), form)
	if err != nil {
		logger.Error().Err(err).Msg("")
		return err
	}
	resp.Body.Close()

	requestURL := fmt.Sprintf("http://%s/api/v1/probe", *addr)
	logger.Info().Msg("Refreshing probes")

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		logger.Error().Err(err).Msg("")
		return err
	}

	resp2, err := client.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("Probe refresh failure")
		return err
	}
	defer resp2.Body.Close()

	var probes []api.Probe
	err = json.NewDecoder(resp2.Body).Decode(&probes)
	if err != nil {
		logger.Error().Err(err).Msg("Probe body read failure")
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

	logger.Debug().Str("policies", fmt.Sprintf("%+v", policies)).Msg("Parsed policies")

	// download policies
	// TODO: compare policies with existing policies and delete those no longer needed
	for name, versions := range policies {
		for _, version := range versions {
			policy_name := fmt.Sprintf("%s_%s", name, version)
			file_name := fmt.Sprintf("bin/%s", policy_name)
			if _, err := os.Stat(file_name); err != nil {
				err = downloadFile(file_name, fmt.Sprintf("http://%s/api/v1/policy/%s/%s", *addr, name, version), client)
				if err != nil {
					logger.Error().Err(err).Str("file", policy_name).Msg("Error downloading policy")
				}
			} else {
				logger.Info().Str("file", policy_name).Msg("Policy already exists, skipping download")
			}
		}
	}

	requestURL = fmt.Sprintf("http://%s/logout", *addr)
	req, err = http.NewRequest("GET", requestURL, nil)
	if err != nil {
		logger.Error().Err(err).Msg("")
		return err
	}

	resp3, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp3.Body.Close()
	return nil
}

func downloadFile(filepath string, url string, client *http.Client) error {

	logger.Info().Str("filepath", filepath).Msg("Downloading")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error().Err(err).Msg("")
		return err
	}

	resp, err := client.Do(req)
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
	logger.Info().Str("filepath", filePath).Msg("Setting execute permissions")
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
