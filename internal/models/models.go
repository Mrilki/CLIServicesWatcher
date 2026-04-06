package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

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

type CheckType string

const (
	CheckTypeTCP  CheckType = "tcp"
	CheckTypeHTTP CheckType = "http"
	CheckTypeDNS  CheckType = "dns"
)

func (c CheckType) IsValid() bool {
	switch c {
	case CheckTypeTCP, CheckTypeHTTP, CheckTypeDNS:
		return true
	default:
		return false
	}

}

type Target struct {
	Name    string    `json:"name"`
	Address string    `json:"address"`
	Timeout *int      `json:"timeout,omitempty"`
	Type    CheckType `json:"type"`
}

func (target *Target) GetTimeoutDuration(globalTimeout time.Duration) time.Duration {
	if target.Timeout == nil || *target.Timeout <= 0 {
		return globalTimeout
	}
	return time.Duration(*target.Timeout) * time.Second
}

type Result struct {
	Name       string    `json:"name"`
	Address    string    `json:"address"`
	StatusCode *int      `json:"status_code,omitempty"`
	Type       CheckType `json:"type"`
	Latency    Duration  `json:"latency"`
	Error      string    `json:"error,omitempty"`
	Success    bool      `json:"success"`
}

func (result *Result) SetLatency(duration time.Duration) {
	result.Latency = Duration(duration)
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
		if target.Address == "" {
			return fmt.Errorf("the target[%d].address is empty", i)
		}
		if !target.Type.IsValid() {
			return fmt.Errorf("the target[%d].type is invalid: %s", i, target.Type)
		}
	}
	return nil
}
func (conf *Config) Normalize() {
	for i := range conf.Targets {
		if conf.Targets[i].Name == "" {
			conf.Targets[i].Name = conf.Targets[i].Address
		}
	}
}

func GetDefaultConf() *Config {
	return &Config{
		Timeout: 10,
		Targets: []Target{
			{Name: "Google", Address: "http://www.google.com", Type: "http"},
			{Name: "Yandex", Address: "http://www.yandex.ru", Type: "http"},
		},
	}
}
