package models

import (
	"errors"
	"fmt"
	"time"
)

type Target struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Result struct {
	Name       string        `json:"name"`
	URL        string        `json:"url"`
	StatusCode int           `json:"status_code"`
	Latency    time.Duration `json:"latency"`
	Error      string        `json:"error,omitempty"`
	Success    bool          `json:"success"`
}

type Config struct {
	Targets []Target `json:"targets"`
	Timeout int      `json:"timeout"`
}

func (conf *Config) GetTimeoutDuration() time.Duration {
	if conf.Timeout <= 0 {
		return 10 * time.Second
	}
	return time.Duration(conf.Timeout) * time.Second
}

func (conf *Config) Validate() error {
	if len(conf.Targets) == 0 {
		return errors.New("the target list is empty")
	}
	for i, target := range conf.Targets {
		if target.Name == "" {
			conf.Targets[i].Name = target.URL
		}
		if target.URL == "" {
			return fmt.Errorf("the target[%d].url is empty", i)
		}
	}
	return nil
}

func (result Result) String() string {
	if result.Success {
		return fmt.Sprintf("Success, %s, %d, %v", result.Name, result.StatusCode, result.Latency)
	}
	return fmt.Sprintf("Error, %s, %d, %v, [%s]", result.Name, result.StatusCode, result.Latency, result.Error)
}

func GetDefaultConf() *Config {
	return &Config{
		Timeout: 10,
		Targets: []Target{
			{Name: "Google", URL: "http://www.google.com"},
			{Name: "Yandex", URL: "http://www.yandex.ru"},
		},
	}
}
