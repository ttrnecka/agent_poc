package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ttrnecka/agent_poc/api"
	"github.com/ttrnecka/agent_poc/ws"
)

var addr = flag.String("addr", "localhost:8888", "http service address")
var source = flag.String("source", "collector1", "name of collector")

var mu sync.Mutex

// reads the collector config and pulls required policies and their versions
func refresh() error {
	requestURL := fmt.Sprintf("http://%s/api/v1/probe", *addr)
	res, err := http.Get(requestURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var probes []api.Probe
	// bodyBytes, err := io.ReadAll(res.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// bodyString := string(bodyBytes)
	// fmt.Println(bodyString)
	err = json.NewDecoder(res.Body).Decode(&probes)
	if err != nil {
		return err
	}

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
	// fmt.Printf("%v\n", policies)

	// download
	for name, versions := range policies {
		for _, version := range versions {
			file_name := fmt.Sprintf("bin/%s_%s", name, version)
			if _, err := os.Stat(file_name); err != nil {
				err = DownloadFile(file_name, fmt.Sprintf("http://%s/api/v1/policy/%s/%s", *addr, name, version))
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
func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			mes := ws.NewMessage(0, "", "", "")

			err := c.ReadJSON(&mes)
			// _, mes, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			messages <- Message{Message: mes, c: c}
		}
	}()

	go func() {
		messageHandler()
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	err = refresh()
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			// err := c.WriteMessage(websocket.TextMessage, []byte("ONLINE"))
			mu.Lock()
			err := c.WriteJSON(ws.NewMessage(ws.MSG_ONLINE, *source, "hub", "Collector is online"))
			mu.Unlock()
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			// err := c.WriteMessage(websocket.TextMessage, []byte("OFFLINE"))
			mu.Lock()
			err := c.WriteJSON(ws.NewMessage(ws.MSG_OFFLINE, *source, "hub", "Collector is going offline"))
			mu.Unlock()
			if err != nil {
				log.Println("write:", err)
				return
			}
			mu.Lock()
			err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			mu.Unlock()
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func DownloadFile(filepath string, url string) error {

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
