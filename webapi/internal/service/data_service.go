package service

import (
	"os"
	"path/filepath"
)

type DataService interface {
	Collectors() ([]string, error)
	CollectorDevices(string) ([]string, error)
	CollectorDeviceEndpoints(string, string) ([]string, error)
	CollectorDeviceEndpointData(string, string, string) ([]byte, error)
}

type dataService struct {
	baseDir string
}

func NewDataService(dir string) DataService {
	return &dataService{baseDir: dir}
}

func (s *dataService) Collectors() ([]string, error) {
	dirs, err := os.ReadDir(s.baseDir)
	if err != nil {
		return nil, err
	}
	collectors := []string{}
	for _, d := range dirs {
		if d.IsDir() {
			collectors = append(collectors, d.Name())
		}
	}
	return collectors, nil
}

func (s *dataService) CollectorDevices(name string) ([]string, error) {
	collectorPath := filepath.Join(s.baseDir, name)
	dirs, err := os.ReadDir(collectorPath)
	if err != nil {
		return nil, err
	}
	devices := []string{}
	for _, d := range dirs {
		if d.IsDir() {
			devices = append(devices, d.Name())
		}
	}
	return devices, nil
}

func (s *dataService) CollectorDeviceEndpoints(collector, device string) ([]string, error) {
	devicePath := filepath.Join(s.baseDir, collector, device)
	entries, err := os.ReadDir(devicePath)
	if err != nil {
		return nil, err
	}

	endpoints := []string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			endpoints = append(endpoints, entry.Name())
		}
	}
	return endpoints, nil
}

func (s *dataService) CollectorDeviceEndpointData(collector, device, endpoint string) ([]byte, error) {
	filePath := filepath.Join(s.baseDir, collector, device, endpoint)
	return os.ReadFile(filePath)
}
