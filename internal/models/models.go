package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Target struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Timeout *int   `json:"timeout,omitempty"`
}

func (target *Target) GetTimeoutDuration(globalTimeout time.Duration) time.Duration {
	if target.Timeout == nil {
		return globalTimeout
	}
	return time.Duration(*target.Timeout) * time.Second
}

type Result struct {
	Name       string   `json:"name"`
	URL        string   `json:"url"`
	StatusCode int      `json:"status_code"`
	Latency    Duration `json:"latency"`
	Error      string   `json:"error,omitempty"`
	Success    bool     `json:"success"`
}

func (result *Result) SetLatency(duration time.Duration) {
	result.Latency = Duration(duration)
}

type Config struct {
	Targets []Target `json:"targets"`
	Timeout int      `json:"timeout"`
}

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Duration(d).String() + `"`), nil
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(dur)
	return nil
}

func (d Duration) String() string {
	return time.Duration(d).String()
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
