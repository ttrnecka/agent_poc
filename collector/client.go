package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/ttrnecka/agent_poc/webapi/api"
)

var httpClient *http.Client

func HTTPClient() *http.Client {
	return httpClient
}

func init() {
	// Initialize an HTTP client. We'll use this for every connection.
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	// New returns an error but currenly it will always be null
	jar, _ := cookiejar.New(nil)

	httpClient = &http.Client{
		Timeout:   time.Minute,
		Transport: t,
		Jar:       jar,
	}
}

func ApiLogin() error {
	form := url.Values{}
	form.Set("username", "admin")
	form.Set("password", "password")

	resp, err := httpClient.PostForm(fmt.Sprintf("http://%s/api/login", *addr), form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// need to read the body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		bodyString := strings.TrimSpace(string(bodyBytes))
		err = fmt.Errorf("login: %s", bodyString)
		return err
	}
	return nil
}

func ApiLogout() error {
	requestURL := fmt.Sprintf("http://%s/api/logout", *addr)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// need to read the body
	defer io.ReadAll(resp.Body)

	// need to read the body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		bodyString := strings.TrimSpace(string(bodyBytes))
		err = fmt.Errorf("logout: %s", bodyString)
		return err
	}
	return nil
}

func ApiGetProbes() ([]api.Probe, error) {
	requestURL := fmt.Sprintf("http://%s/api/v1/probe", *addr)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var probes []api.Probe
	err = json.NewDecoder(resp.Body).Decode(&probes)
	if err != nil {
		return nil, err
	}
	return probes, nil
}
